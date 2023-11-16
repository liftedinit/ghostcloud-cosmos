package types

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

const (
	InvalidCreatorAddress          = "invalid creator address: %s"
	InvalidDomain                  = "invalid domain: %v"
	InvalidName                    = "invalid name: %v"
	InvalidDNSDomain               = "invalid DNS domain: %s"
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

func validateDomain(domain string) error {
	if domain != "" && !govalidator.IsDNSName(domain) {
		return fmt.Errorf(InvalidDNSDomain, domain)
	}
	return nil
}

func validateName(name string) error {
	if name == "" {
		return fmt.Errorf(NameShouldNotBeEmpty)
	}
	if govalidator.HasWhitespace(name) {
		return fmt.Errorf(NameShouldNotContainWhitespace, name)
	}
	if !govalidator.IsASCII(name) {
		return fmt.Errorf(NameShouldContainASCII, name)
	}
	return nil
}
