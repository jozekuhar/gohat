package routes

const (
	PathOrganizationSlug = "organizationSlug"
	PathMembershipID     = "membershipID"
	PathInvitationToken  = "invitationToken"
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
	HXOrgMembershipsUpdate       = "/%s/memberships/%s/update"
	HXOrgMembershipsUpdatePath   = "/{" + PathOrganizationSlug + "}/memberships/{" + PathMembershipID + "}/update"
	HXOrgInvitationsCreate       = "/hx/%s/invitations/create"
	HXOrgInvitationsCreatePath   = "/hx/{" + PathOrganizationSlug + "}/invitations/create"
	InvitationsJoin              = "/invitations/%s/join"
	InvitationsJoinPath          = "/invitations/{" + PathInvitationToken + "}/join"
	HXInvitationsAccept          = "/invitations/%s/accept"
	HXInvitationsAcceptPath      = "/invitations/{" + PathInvitationToken + "}/accept"
	HXInvitationsDecline         = "/invitations/%s/decline"
	HXInvitationsDeclinePath     = "/invitations/{" + PathInvitationToken + "}/decline"

	OrgChannels             = "/%s/channels"
	OrgChannelsPath         = "/{" + PathOrganizationSlug + "}/channels"
	HXOrgChannelsCreate     = "/%s/channels/create"
	HXOrgChannelsCreatePath = "/{" + PathOrganizationSlug + "}/channels/create"
)
