package cloudobject

func Exists(obj CloudObject) (bool, error) {
	if err := obj.Read(); err != nil {
		if IsNotExistsError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil

}
