package routes

const (
	PathOrganizationSlug = "organizationSlug"
)

const (
	Static = "/static/*"
)

const (
	Index                = "/"
	IExampleModal        = "/i/example/modal"
	IExampleToast        = "/i/example/toast"
	IExampleLongRequest  = "/i/example/long-request"
	IExampleSSERequest   = "/i/example/sse-request"
	Chat                 = "/chat"
	WSChat               = "/ws/chat"
	Login                = "/login"
	ILoginGoogle         = "/i/login/google"
	ILoginGoogleCallback = "/i/login/google/callback"
	ILogout              = "/i/logout"

	// Organization
	HXOrganizationCreate = "/hx/organization/create"

	Dashboard     = "/%s"
	DashboardPath = "/{" + PathOrganizationSlug + "}"

	Memberships           = "/%s/memberships"
	MembershipsPath       = "/{" + PathOrganizationSlug + "}/memberships"
	MembershipsCreate     = "/%s/memberships/create"
	MembershipsCreatePath = "/{" + PathOrganizationSlug + "}/memberships/create"

	Channels     = "/%s/channels"
	ChannelsPath = "/{" + PathOrganizationSlug + "}/channels"
)
