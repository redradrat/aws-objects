package cloudobject

// NotExistsError is returned when a Cloud Object does not exist
type NotExistsError struct {
	Message string
}

func (e NotExistsError) Error() string {
	return e.Message
}

func IsNotExistsError(err error) bool {
	_, ok := err.(NotExistsError)
	return ok
}

func IgnoreNotExistsError(err error) error {
	if IsNotExistsError(err) {
		return nil
	}
	return err
}

// NotReadyError is returned when a Cloud Object does not exist
type NotReadyError struct {
	Message string
}

func (e NotReadyError) Error() string {
	return e.Message
}

func IsNotReadyError(err error) bool {
	_, ok := err.(NotReadyError)
	return ok
}

func IgnoreNotReadyError(err error) error {
	if IsNotReadyError(err) {
		return nil
	}
	return err
}

// NotExistsError is returned when a Cloud Object does not exist
type AmbiguousIdentifierError struct {
	Message string
}

func (e AmbiguousIdentifierError) Error() string {
	return e.Message
}

func IsAmbiguousIdentifierError(err error) bool {
	_, ok := err.(AmbiguousIdentifierError)
	return ok
}

func IgnoreAmbiguousIdentifierError(err error) error {
	if IsAmbiguousIdentifierError(err) {
		return nil
	}
	return err
}

// AlreadyExistsError is returned when a Cloud Object already exists
type AlreadyExistsError struct {
	Message string
}

func (e AlreadyExistsError) Error() string {
	return e.Message
}

func IsAlreadyExistsError(err error) bool {
	_, ok := err.(AlreadyExistsError)
	return ok
}

func IgnoreAlreadyExistsError(err error) error {
	if IsAlreadyExistsError(err) {
		return nil
	}
	return err
}

// SpecInvalidError is returned when a CloudObjectSpec is invalid for the current action
type SpecInvalidError struct {
	Message string
}

func (e SpecInvalidError) Error() string {
	return e.Message
}

func IsCloudSpecInvalidError(err error) bool {
	_, ok := err.(SpecInvalidError)
	return ok
}

func IgnoreCloudSpecInvalidError(err error) error {
	if IsCloudSpecInvalidError(err) {
		return nil
	}
	return err
}

// OptsInvalidError is returned when a an options object (e.g. DeleteOpts) is invalid for the current action
type OptsInvalidError struct {
	Message string
}

func (e OptsInvalidError) Error() string {
	return e.Message
}

func IsOptsInvalidError(err error) bool {
	_, ok := err.(OptsInvalidError)
	return ok
}

func IgnoreOptsInvalidError(err error) error {
	if IsOptsInvalidError(err) {
		return nil
	}
	return err
}

// IdCollisionError is returned when a CloudObject cannot be created due to collisions with existing objects
type IdCollisionError struct {
	Message string
}

func (e IdCollisionError) Error() string {
	return e.Message
}

func IsIdCollisionError(err error) bool {
	_, ok := err.(IdCollisionError)
	return ok
}

func IgnoreIdCollisionError(err error) error {
	if IsIdCollisionError(err) {
		return nil
	}
	return err
}
