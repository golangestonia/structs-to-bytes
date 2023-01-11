package main

import (
	"fmt"
	"reflect"

	"github.com/golang-estonia/structs-to-bytes/est"
	"github.com/zeebo/errs"
)

func Encode(v any) (data []byte, err error) {
	return encode(reflect.ValueOf(v))
}

func encode(v reflect.Value) (data []byte, err error) {
	// WARNING: THIS CODE MAY CONTAIN SUBTLE BUGS
	// CODE HAS NOT BEEN PROPERLY REVIEWED AND TESTED!

	kind := v.Kind()
	if kind == reflect.Pointer {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Uint64:
		data = est.AppendUint64(data, v.Uint())
		return data, nil
	case reflect.String:
		data = est.AppendString(data, v.String())
		return data, nil
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			sub, err := encode(f)
			if err != nil {
				return nil, errs.Wrap(err)
			}

			// TODO: handle pointer to struct here
			if f.Kind() == reflect.Struct {
				data = est.AppendBytes(data, sub)
			} else {
				data = append(data, sub...)
			}
		}
		return data, nil
	case reflect.Slice:
		data, err = est.AppendMessage(data, func() (data []byte, err error) {
			n := v.Len()
			data = est.AppendUint64(data, uint64(n))
			for i := 0; i < n; i++ {
				f := v.Index(i)
				sub, err := encode(f)
				if err != nil {
					return nil, errs.Wrap(err)
				}

				if f.Kind() == reflect.Struct {
					data = est.AppendBytes(data, sub)
				} else {
					data = append(data, sub...)
				}
			}
			return data, nil
		})
		return data, errs.Wrap(err)
	default:
		return nil, fmt.Errorf("unhandled %v", v.Kind())
	}

	return data, nil
}

func Decode(data []byte, v any) (err error) {
	_, err = decode(data, reflect.ValueOf(v))
	return errs.Wrap(err)
}

func decode(data []byte, v reflect.Value) (rest []byte, err error) {
	// WARNING: THIS CODE MAY CONTAIN SUBTLE BUGS
	// CODE HAS NOT BEEN PROPERLY REVIEWED AND TESTED!

	kind := v.Kind()
	if kind == reflect.Pointer {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Uint64:
		var x uint64
		x, data, err = est.ReadUint64(data)
		if err != nil {
			return nil, err
		}
		v.SetUint(x)
		return data, nil

	case reflect.String:
		var x string
		x, data, err = est.ReadString(data)
		if err != nil {
			return nil, err
		}
		v.SetString(x)
		return data, nil

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			// TODO: handle pointer to struct here
			if f.Kind() == reflect.Struct {
				data, err = est.ReadMessage(data, func(msgdata []byte) (err error) {
					msgdata, err = decode(msgdata, f)
					if err != nil {
						return errs.Wrap(err)
					}
					if len(msgdata) != 0 {
						return errs.New("too much data")
					}
					return nil
				})
				if err != nil {
					return nil, errs.Wrap(err)
				}
			} else {
				data, err = decode(data, f)
				if err != nil {
					return nil, errs.Wrap(err)
				}
			}
		}
		return data, nil
	case reflect.Slice:
		return est.ReadMessage(data, func(data []byte) error {
			var n uint64
			n, data, err = est.ReadUint64(data)
			if err != nil {
				return err
			}
			v.Set(reflect.MakeSlice(v.Type(), int(n), int(n)))

			for i := 0; i < int(n); i++ {
				f := v.Index(i)
				if f.Kind() == reflect.Struct {
					data, err = est.ReadMessage(data, func(msgdata []byte) (err error) {
						msgdata, err = decode(msgdata, f)
						if err != nil {
							return errs.Wrap(err)
						}
						if len(msgdata) != 0 {
							return errs.New("too much data")
						}
						return nil
					})
					if err != nil {
						return errs.Wrap(err)
					}
				} else {
					data, err = decode(data, f)
					if err != nil {
						return errs.Wrap(err)
					}
				}
			}

			if len(data) != 0 {
				return errs.New("too much data")
			}

			return nil
		})

	default:
		return data, fmt.Errorf("unhandled %v", v.Kind())
	}

	return data, nil
}
