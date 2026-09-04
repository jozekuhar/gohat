package routes

const (
	PathOrganizationSlug = "organizationSlug"
	PathInvitationID     = "invitationID"
)

const (
	Static = "/static/*"
)

const (
	Root                         = "/"
	Login                        = "/login"
	HXLogin                      = "/hx/login"
	Register                     = "/register"
	HXRegister                   = "/hx/register"
	HXSignInGoogle               = "/hx/signin/google"
	CallbackSignInGoogle         = "/callback/signin/google"
	HXLogout                     = "/hx/logout"
	PrivacyPolicy                = "/privacy"
	TermsOfService               = "/terms"
	HXSidebarOrganizations       = "/hx/sidebar/organizations"
	HXSidebarOrganizationsCreate = "/hx/sidebar/organizations/create"
	OrgRoot                      = "/%s"
	OrgRootPath                  = "/{" + PathOrganizationSlug + "}"
	OrgDashboard                 = "/%s/dashboard"
	OrgDashboardPath             = "/{" + PathOrganizationSlug + "}/dashboard"
	OrgMemberships               = "/%s/memberships"
	OrgMembershipsPath           = "/{" + PathOrganizationSlug + "}/memberships"
	HXOrgMembershipsCreate       = "/hx/%s/memberships/create"
	HXOrgMembershipsCreatePath   = "/hx/{" + PathOrganizationSlug + "}/memberships/create"

	// TO DO BELLOW
	Invitations           = "/%s/invitations"
	InvitationsPath       = "/{" + PathOrganizationSlug + "}/invitations"
	InvitationsDetail     = "/%s/invitations/%s"
	InvitationsDetailPath = "/{" + PathOrganizationSlug + "}/invitations/{" + PathInvitationID + "}"
	InvitationsJoin       = "/%s/invitations/join"
	InvitationsJoinPath   = "/{" + PathOrganizationSlug + "}/invitations/join"

	Channels     = "/%s/channels"
	ChannelsPath = "/{" + PathOrganizationSlug + "}/channels"
)
