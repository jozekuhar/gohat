package routes

const (
	PathOrganizationSlug = "organizationSlug"
	PathInvitationID     = "invitationID"
)

const (
	Static = "/static/*"
)

const (
	Index = "/"

	Login                = "/login"
	HXLogin              = "/hx/login"
	Register             = "/register"
	HXRegister           = "/hx/register"
	HXSignInGoogle       = "/hx/signin/google"
	CallbackSignInGoogle = "/callback/signin/google"
	HXLogout             = "/hx/logout"

	// TO DO BELLOW
	HXOrganizationCreate = "/hx/organization/create"

	Dashboard     = "/%s"
	DashboardPath = "/{" + PathOrganizationSlug + "}"

	Memberships           = "/%s/memberships"
	MembershipsPath       = "/{" + PathOrganizationSlug + "}/memberships"
	Invitations           = "/%s/invitations"
	InvitationsPath       = "/{" + PathOrganizationSlug + "}/invitations"
	InvitationsDetail     = "/%s/invitations/%s"
	InvitationsDetailPath = "/{" + PathOrganizationSlug + "}/invitations/{" + PathInvitationID + "}"
	InvitationsJoin       = "/%s/invitations/join"
	InvitationsJoinPath   = "/{" + PathOrganizationSlug + "}/invitations/join"

	Channels     = "/%s/channels"
	ChannelsPath = "/{" + PathOrganizationSlug + "}/channels"
)
