package addrparam

//func convertHostParams(orig HostNameParameters) *hostNameParameters {
//	if params, ok := orig.(*hostNameParameters); ok {
//		return params
//	}
//
//	paramsBuilder := HostNameParametersBuilder{}
//	return paramsBuilder.
//		// general settings
//		AllowIPAddress(orig.AllowsIPAddress()).
//		AllowBracketedIPv6(orig.AllowsBracketedIPv6()).
//		AllowBracketedIPv4(orig.AllowsBracketedIPv4()).
//		SetEmptyLoopback(orig.EmptyIsLoopback()).
//		AllowPort(orig.AllowsPort()).
//		AllowService(orig.AllowsService()).
//		ExpectPort(orig.ExpectsPort()).
//		AllowEmpty(orig.AllowsEmpty()).
//		NormalizeToLowercase(orig.NormalizesToLowercase()).
//		SetIPAddressParameters(orig.GetIPAddressParameters()).
//		//
//		ToParams().(*hostNameParameters)
//}

func CopyHostNameParams(orig HostNameParameters) HostNameParameters {
	if p, ok := orig.(*hostNameParameters); ok {
		return p
	}
	return new(HostNameParametersBuilder).Set(orig).ToParams()
}

//func DefaultHostNameParams() HostNameParameters {
//	xxx use builder instead xxx
//	return &hostNameParameters{}
//}

type HostNameParameters interface {
	// AllowsEmpty determines if an empty host string, when not a valid address, is considered valid.
	// The parser will first parse as an empty address, if allowed by the nested IPAddressStringParameters.
	// Otherwise, it will be considered an empty host if this returns true, or an invalid host if it returns false.
	AllowsEmpty() bool

	//xxxx gotta defer to address on this xxx
	//EmptyStrParsedAs() EmptyStrOption

	// Indicates the version to prefer when resolving host names.
	GetPreferredVersion() IPVersion

	//EmptyIsLoopback() bool
	AllowsBracketedIPv4() bool
	AllowsBracketedIPv6() bool
	NormalizesToLowercase() bool
	AllowsIPAddress() bool
	AllowsPort() bool
	AllowsService() bool
	ExpectsPort() bool

	GetIPAddressParameters() IPAddressStringParameters

	//ToAddressOptionsBuilder() IPAddressStringParametersBuilder
}

// hostNameParameters has parameters for parsing IP address strings
// They are immutable and can be constructed using an HostNameParametersBuilder
type hostNameParameters struct {
	ipParams ipAddressStringParameters

	//emptyStringOption EmptyStrOption

	preferredVersion IPVersion

	noEmpty, noBracketedIPv4, noBracketedIPv6,
	noNormalizeToLower, noIPAddress, noPort, noService, expectPort bool
}

//func (params *hostNameParameters) ToAddressOptionsBuilder() IPAddressStringParametersBuilder {xxx no longer use this xxxx
//	return params.ipParams.ToBuilder()
//}

func (params *hostNameParameters) AllowsEmpty() bool {
	return !params.noEmpty
}

//func (params *hostNameParameters) EmptyStrParsedAs() EmptyStrOption {
//	return params.emptyStringOption
//}

func (params *hostNameParameters) GetPreferredVersion() IPVersion {
	return params.preferredVersion
}

func (params *hostNameParameters) AllowsBracketedIPv4() bool {
	return !params.noBracketedIPv4
}

func (params *hostNameParameters) AllowsBracketedIPv6() bool {
	return !params.noBracketedIPv6
}

func (params *hostNameParameters) NormalizesToLowercase() bool {
	return !params.noNormalizeToLower
}

func (params *hostNameParameters) AllowsIPAddress() bool {
	return !params.noIPAddress
}

func (params *hostNameParameters) AllowsPort() bool {
	return !params.noPort
}

func (params *hostNameParameters) AllowsService() bool {
	return !params.noService
}

func (params *hostNameParameters) ExpectsPort() bool {
	return params.expectPort
}

func (params *hostNameParameters) GetIPAddressParameters() IPAddressStringParameters {
	return &params.ipParams
}

// HostNameParametersBuilder builds a hostNameParameters
type HostNameParametersBuilder struct {
	hostNameParameters

	ipAddressBuilder IPAddressStringParametersBuilder
}

//func (params *hostNameParameters) ToAddressOptionsBuilder() IPAddressStringParametersBuilder {xxx no longer use this xxxx
//	return params.ipParams.ToBuilder()
//}

//func ToIPAddressParametersBuilder(params HostNameParameters) *IPAddressStringParametersBuilder {
//	return ToIPAddressStringParamsBuilder(params.GetIPAddressParameters())
//}

//func ToHostNameParametersBuilder(params HostNameParameters) *HostNameParametersBuilder {
//	var result HostNameParametersBuilder
//	if p, ok := params.(*hostNameParameters); ok {
//		result.hostNameParameters = *p
//	} else {
//		result.hostNameParameters = hostNameParameters{
//			emptyStringOption:  params.EmptyStrParsedAs(),
//			preferredVersion:   params.GetPreferredVersion(),
//			noEmpty:            !params.AllowsEmpty(),
//			noBracketedIPv4:    !params.AllowsBracketedIPv4(),
//			noBracketedIPv6:    !params.AllowsBracketedIPv6(),
//			noNormalizeToLower: !params.NormalizesToLowercase(),
//			noIPAddress:        !params.AllowsIPAddress(),
//			noPort:             !params.AllowsPort(),
//			noService:          !params.AllowsService(),
//			expectPort:         params.ExpectsPort(),
//		}
//	}
//	result.SetIPAddressParameters(params.GetIPAddressParameters())
//	return &result
//}

func (builder *HostNameParametersBuilder) ToParams() HostNameParameters {
	// We do not return a pointer to builder.params because that would make it possible to change a ipAddressStringParameters
	// by continuing to use the same builder,
	// and we want immutable objects for thread-safety,
	// so we cannot allow it
	result := builder.hostNameParameters
	result.ipParams = *builder.ipAddressBuilder.ToParams().(*ipAddressStringParameters)
	return &result
}

func (builder *HostNameParametersBuilder) GetIPAddressParametersBuilder() (result *IPAddressStringParametersBuilder) {
	result = &builder.ipAddressBuilder
	result.parent = builder
	return
}

func (builder *HostNameParametersBuilder) Set(params HostNameParameters) *HostNameParametersBuilder {
	//var result HostNameParametersBuilder
	if p, ok := params.(*hostNameParameters); ok {
		builder.hostNameParameters = *p
	} else {
		builder.hostNameParameters = hostNameParameters{
			//emptyStringOption:  params.EmptyStrParsedAs(),
			preferredVersion:   params.GetPreferredVersion(),
			noEmpty:            !params.AllowsEmpty(),
			noBracketedIPv4:    !params.AllowsBracketedIPv4(),
			noBracketedIPv6:    !params.AllowsBracketedIPv6(),
			noNormalizeToLower: !params.NormalizesToLowercase(),
			noIPAddress:        !params.AllowsIPAddress(),
			noPort:             !params.AllowsPort(),
			noService:          !params.AllowsService(),
			expectPort:         params.ExpectsPort(),
		}
	}
	builder.SetIPAddressParameters(params.GetIPAddressParameters())
	return builder
}

func (builder *HostNameParametersBuilder) SetIPAddressParameters(params IPAddressStringParameters) *HostNameParametersBuilder {
	//builder.ipAddressBuilder = *ToIPAddressStringParamsBuilder(params)
	builder.ipAddressBuilder.Set(params)
	return builder
}

func (builder *HostNameParametersBuilder) AllowEmpty(allow bool) *HostNameParametersBuilder {
	builder.hostNameParameters.noEmpty = !allow
	return builder
}

//func (builder *HostNameParametersBuilder) SetEmptyLoopback(isLoopback bool) *HostNameParametersBuilder {
//	builder.hostNameParameters.emptyIsNotLoopback = !isLoopback
//	return builder
//}
//
//func (builder *HostNameParametersBuilder) ParseEmptyStrAs(option EmptyStrOption) *HostNameParametersBuilder {
//	builder.hostNameParameters.emptyStringOption = option
//	if option != NoAddressOption {
//		builder.AllowEmpty(true)
//	}
//	return builder
//}

func (builder *HostNameParametersBuilder) SetPreferredVersion(version IPVersion) *HostNameParametersBuilder {
	builder.hostNameParameters.preferredVersion = version
	return builder
}

func (builder *HostNameParametersBuilder) AllowBracketedIPv4(allow bool) *HostNameParametersBuilder {
	builder.hostNameParameters.noBracketedIPv4 = !allow
	return builder
}

func (builder *HostNameParametersBuilder) AllowBracketedIPv6(allow bool) *HostNameParametersBuilder {
	builder.hostNameParameters.noBracketedIPv6 = !allow
	return builder
}

func (builder *HostNameParametersBuilder) NormalizeToLowercase(allow bool) *HostNameParametersBuilder {
	builder.hostNameParameters.noNormalizeToLower = !allow
	return builder
}

func (builder *HostNameParametersBuilder) AllowIPAddress(allow bool) *HostNameParametersBuilder {
	builder.hostNameParameters.noIPAddress = !allow
	return builder
}

func (builder *HostNameParametersBuilder) AllowPort(allow bool) *HostNameParametersBuilder {
	builder.hostNameParameters.noPort = !allow
	return builder
}

func (builder *HostNameParametersBuilder) AllowService(allow bool) *HostNameParametersBuilder {
	builder.hostNameParameters.noService = !allow
	return builder
}

func (builder *HostNameParametersBuilder) ExpectPort(expect bool) *HostNameParametersBuilder {
	builder.hostNameParameters.expectPort = expect
	return builder
}