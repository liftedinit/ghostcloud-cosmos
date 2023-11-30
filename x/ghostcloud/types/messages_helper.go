package types

const (
	InvalidCreatorAddress          = "invalid creator address: %s"
	InvalidDomain                  = "invalid domain: %v"
	CreatorShouldNotBeEmpty        = "creator should not be empty"
	NameShouldNotBeEmpty           = "name should not be empty"
	NameShouldNotContainWhitespace = "name should not contain whitespace: %s"
	NameShouldContainASCII         = "name should contain ascii characters only: %s"
	NameTooLong                    = "name is too long: %s"
	DescriptionTooLong             = "description is too long: %s"
	MetaIsRequired                 = "meta is required"
	PayloadTooBig                  = "payload is too big: %d > %d"
	PayloadIsRequired              = "payload is required"
	UncompressedSizeTooBig         = "total uncompressed size is too big: %d > %d"
	IndexHtmlNotFound              = "index.html not found"
	NothingToUpdate                = "nothing to update"
)
