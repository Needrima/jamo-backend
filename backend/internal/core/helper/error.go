package helper

import "errors"

var NEWSLETTER_MAIL_ERROR = errors.New("could not send mail on newsletter subscription")
var CONTACT_MAIL_ERROR = errors.New("something went wrong, please try later")
var USER_ALREADY_A_SUBSCRIBER = errors.New("already subscribed to newsletter")
var INVALID_PAYLOAD = errors.New("invalid field in payload body")
