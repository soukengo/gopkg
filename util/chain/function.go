package chain

// HandleWithResultFunc functions to be executed
// ret: function return
// cont: if execute the next function
// err: executed error
type HandleWithResultFunc[T any] func() (ret T, cont bool, err error)

// HandleFunc functions to be executed
// err: executed error, when error is returned, chain will be stopped
type HandleFunc func() (err error)

// CallWithResult makes a call chain of many functions with result
func CallWithResult[T any](chains ...HandleWithResultFunc[T]) (ret T, err error) {
	for _, chain := range chains {
		var cont bool
		ret, cont, err = chain()
		if !cont {
			break
		}
	}
	return
}

// Call makes a call chain of many functions without result
func Call(chains ...HandleFunc) error {
	for _, chain := range chains {
		if chain == nil {
			continue
		}
		if err := chain(); err != nil {
			return err
		}
	}
	return nil
}
