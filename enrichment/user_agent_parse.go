package enrichment

import (
	"github.com/jitsucom/eventnative/appconfig"
	"github.com/jitsucom/eventnative/jsonutils"
	"github.com/jitsucom/eventnative/logging"
	"github.com/jitsucom/eventnative/parsers"
	"github.com/jitsucom/eventnative/useragent"
)

const UserAgentParse = "user_agent_parse"

type UserAgentParseRule struct {
	source                  *jsonutils.JsonPath
	destination             *jsonutils.JsonPath
	uaResolver              useragent.Resolver
	enrichmentConditionFunc func(map[string]interface{}) bool
}

func NewUserAgentParseRule(source, destination *jsonutils.JsonPath) (*UserAgentParseRule, error) {
	return &UserAgentParseRule{
		source:      source,
		destination: destination,
		uaResolver:  appconfig.Instance.UaResolver,
		//always do enrichment
		enrichmentConditionFunc: func(m map[string]interface{}) bool {
			return true
		},
	}, nil
}

func (uap *UserAgentParseRule) Execute(event map[string]interface{}) {
	if !uap.enrichmentConditionFunc(event) {
		return
	}

	uaIface, ok := uap.source.Get(event)
	if !ok {
		return
	}

	ua, ok := uaIface.(string)
	if !ok {
		return
	}

	parsedUa := uap.uaResolver.Resolve(ua)

	//convert all structs to map[string]interface{} for inner typecasting
	result, err := parsers.ParseInterface(parsedUa)
	if err != nil {
		logging.SystemErrorf("Error converting ua parse node: %v", err)
		return
	}

	err = uap.destination.Set(event, result)
	if err != nil {
		logging.SystemErrorf("Resolved useragent data wasn't set: %v", err)
	}
}

func (uap *UserAgentParseRule) Name() string {
	return UserAgentParse
}
