-- name: CreateOrganization :one
INSERT INTO organizations (id, name, slug)
VALUES (@id, @name, @slug)
RETURNING *;


-- name: ListActiveOrganizations :many
SELECT o.*
FROM organizations AS o
INNER JOIN memberships AS m ON o.id = m.organization_id
WHERE m.user_id = @user_id
  AND m.status = 'active';


-- name: CreateMembership :one
INSERT INTO memberships (id, organization_id, user_id, first_name, last_name, role, permissions, status)
VALUES (@id, @organization_id, @user_id, @first_name, @last_name, @role, @permissions, @status)
RETURNING *;

-- name: ListMemberships :many
SELECT sqlc.embed(m), sqlc.embed(u)
FROM memberships AS m
JOIN users AS u ON m.user_id = u.id
WHERE m.organization_id = @organization_id;


-- name: GetMembership :one
SELECT sqlc.embed(m), sqlc.embed(u)
FROM memberships AS m
JOIN users AS u ON m.user_id = u.id
WHERE m.organization_id = @organization_id
  AND m.id = @membership_id;


-- name: GetActiveMembershipByUserID :one
SELECT sqlc.embed(o), sqlc.embed(m), sqlc.embed(u)
FROM organizations AS o
INNER JOIN memberships AS m ON o.id = m.organization_id
INNER JOIN users AS u ON m.user_id = u.id
WHERE o.slug = @slug
  AND m.user_id = @user_id
  AND m.status = 'active';


-- name: CheckMembershipByEmail :one
SELECT EXISTS (
    SELECT 1 
    FROM memberships AS m
    JOIN users AS u ON m.user_id = u.id 
    WHERE m.organization_id = @organization_id
      AND u.email = @email
);


-- name: UpdateMembership :one
UPDATE memberships
SET first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    updated_at = NOW()
WHERE organization_id = @organization_id
  AND id = @membership_id
RETURNING *;

-- name: UpdateMembershipCanceled :exec
UPDATE memberships
SET status = 'canceled',
    canceled_at = NOW(),
    canceled_by_id = @canceled_by_id
WHERE organization_id = @organization_id
  AND id = @membership_id
  AND status = 'active';

-- name: CreateInvitation :one
INSERT INTO invitations (id, organization_id, inviter_id, email, first_name, last_name, role, permissions, token_hash, expires_at)
VALUES (@id, @organization_id, @inviter_id, @email, @first_name, @last_name, @role, @permissions, @token_hash, @expires_at)
RETURNING *;


-- name: ListInvitations :many
SELECT *
FROM invitations
WHERE organization_id = @organization_id
ORDER BY created_at DESC;


-- name: GetInvitationByTokenHash :one
SELECT sqlc.embed(i), sqlc.embed(o), sqlc.embed(u), sqlc.embed(m)
FROM invitations AS i
INNER JOIN organizations AS o ON i.organization_id = o.id
INNER JOIN users AS u ON i.inviter_id = u.id
INNER JOIN memberships AS m ON i.inviter_id = m.user_id AND i.organization_id = m.organization_id
WHERE i.token_hash = @token_hash;


-- name: GetLatestInvitationByEmail :one
SELECT *
FROM invitations
WHERE organization_id = @organization_id
  AND email = @email
ORDER BY created_at DESC
LIMIT 1;


-- name: UpdateInvitationAcceptedAt :exec
UPDATE invitations
SET accepted_at = NOW()
WHERE id = @invitation_id
  AND accepted_at IS NULL
  AND declined_at IS NULL
  AND canceled_at IS NULL;


-- name: UpdateInvitationDeclinedAt :exec
UPDATE invitations
SET declined_at = NOW()
WHERE id = @invitation_id
  AND accepted_at IS NULL
  AND declined_at IS NULL
  AND canceled_at IS NULL;


-- name: UpdateInvitationCanceled :exec
UPDATE invitations
SET canceled_at = NOW(),
    canceled_by_id = @canceled_by_id
WHERE organization_id = @organization_id
  AND id = @invitation_id
  AND accepted_at IS NULL
  AND declined_at IS NULL
  AND canceled_at IS NULL;

-- name: DeleteInvitation :exec
DELETE FROM invitations
WHERE organization_id = @organization_id
  AND id = @invitation_id;
