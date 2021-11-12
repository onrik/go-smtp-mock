package smtpmock

import "errors"

// HELO command handler
type handlerHelo struct {
	*handler
}

// HELO command handler builder. Returns pointer to new handlerHelo structure
func newHandlerHelo(session sessionInterface, message *message, configuration *configuration) *handlerHelo {
	return &handlerHelo{&handler{session: session, message: message, configuration: configuration}}
}

// HELO handler methods

// Main HELO handler runner
func (handler *handlerHelo) run(request string) {
	handler.clearError()
	handler.clearMessage()

	if handler.isInvalidRequest(request) {
		return
	}

	handler.writeResult(true, request, handler.configuration.msgHeloReceived)
}

// Erases message data
func (handler *handlerHelo) clearMessage() {
	*handler.message = *zeroMessage
}

// Writes handled HELO result to session, message. Always returns true
func (handler *handlerHelo) writeResult(isSuccessful bool, request, response string) bool {
	session, message := handler.session, handler.message
	if !isSuccessful {
		session.addError(errors.New(response))
	}

	message.heloRequest, message.heloResponse, message.helo = request, response, isSuccessful
	session.writeResponse(response)
	return true
}

// Invalid HELO command argument predicate. Returns true and writes result for case when HELO command
// argument is invalid, otherwise returns false
func (handler *handlerHelo) isInvalidCmdArg(request string) bool {
	if !matchRegex(request, ValidHeloComplexCmdRegexPattern) {
		return handler.writeResult(false, request, handler.configuration.msgInvalidCmdHeloArg)
	}

	return false
}

// Returns domain from HELO request
func (handler *handlerHelo) heloDomain(request string) string {
	return regexCaptureGroup(request, ValidHeloComplexCmdRegexPattern, 2)
}

// Custom behaviour for HELO domain. Returns true and writes result for case when HELO domain
// is included in configuration.blacklistedHeloDomains slice
func (handler *handlerHelo) isBlacklistedDomain(request string) bool {
	configuration := handler.configuration
	if !isIncluded(configuration.blacklistedHeloDomains, handler.heloDomain(request)) {
		return false
	}

	return handler.writeResult(false, request, configuration.msgHeloBlacklistedDomain)
}

// Invalid HELO command request complex predicate. Returns true for case when one
// of the chain checks returns true, otherwise returns false
func (handler *handlerHelo) isInvalidRequest(request string) bool {
	return handler.isInvalidCmdArg(request) || handler.isBlacklistedDomain(request)
}
