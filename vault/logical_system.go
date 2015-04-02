package vault

import (
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func NewSystemBackend(core *Core) logical.Backend {
	b := &SystemBackend{Core: core}

	return &framework.Backend{
		PathsSpecial: &logical.Paths{
			Root: []string{
				"mounts/*",
				"auth/*",
				"remount",
				"revoke-prefix/*",
				"policy",
				"policy/*",
				"audit",
				"audit/*",
				"seal", // Must be set for Core.Seal() logic
				"raw/*",
			},
		},

		Paths: []*framework.Path{
			&framework.Path{
				Pattern: "mounts$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.handleMountTable,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["mounts"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["mounts"][1]),
			},

			&framework.Path{
				Pattern: "mounts/(?P<path>.+)",

				Fields: map[string]*framework.FieldSchema{
					"path": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["mount_path"][0]),
					},

					"type": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["mount_type"][0]),
					},
					"description": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["mount_desc"][0]),
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.WriteOperation:  b.handleMount,
					logical.DeleteOperation: b.handleUnmount,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["mount"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["mount"][1]),
			},

			&framework.Path{
				Pattern: "remount",

				Fields: map[string]*framework.FieldSchema{
					"from": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["remount_from"][0]),
					},
					"to": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["remount_to"][0]),
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.WriteOperation: b.handleRemount,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["remount"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["remount"][1]),
			},

			&framework.Path{
				Pattern: "renew/(?P<vault_id>.+)",

				Fields: map[string]*framework.FieldSchema{
					"vault_id": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["vault_id"][0]),
					},
					"increment": &framework.FieldSchema{
						Type:        framework.TypeInt,
						Description: strings.TrimSpace(sysHelp["increment"][0]),
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.WriteOperation: b.handleRenew,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["renew"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["renew"][1]),
			},

			&framework.Path{
				Pattern: "revoke/(?P<vault_id>.+)",

				Fields: map[string]*framework.FieldSchema{
					"vault_id": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["vault_id"][0]),
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.WriteOperation: b.handleRevoke,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["revoke"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["revoke"][1]),
			},

			&framework.Path{
				Pattern: "revoke-prefix/(?P<prefix>.+)",

				Fields: map[string]*framework.FieldSchema{
					"prefix": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["revoke-prefix-path"][0]),
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.WriteOperation: b.handleRevokePrefix,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["revoke-prefix"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["revoke-prefix"][1]),
			},

			&framework.Path{
				Pattern: "auth$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.handleAuthTable,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["auth-table"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["auth-table"][1]),
			},

			&framework.Path{
				Pattern: "auth/(?P<path>.+)",

				Fields: map[string]*framework.FieldSchema{
					"path": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["auth_path"][0]),
					},
					"type": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["auth_type"][0]),
					},
					"description": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["auth_desc"][0]),
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.WriteOperation:  b.handleEnableAuth,
					logical.DeleteOperation: b.handleDisableAuth,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["auth"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["auth"][1]),
			},

			&framework.Path{
				Pattern: "policy$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.handlePolicyList,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["policy-list"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["policy-list"][1]),
			},

			&framework.Path{
				Pattern: "policy/(?P<name>.+)",

				Fields: map[string]*framework.FieldSchema{
					"name": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["policy-name"][0]),
					},
					"rules": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["policy-rules"][0]),
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation:   b.handlePolicyRead,
					logical.WriteOperation:  b.handlePolicySet,
					logical.DeleteOperation: b.handlePolicyDelete,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["policy"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["policy"][1]),
			},

			&framework.Path{
				Pattern: "audit$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.handleAuditTable,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["audit-table"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["audit-table"][1]),
			},

			&framework.Path{
				Pattern: "audit/(?P<path>.+)",

				Fields: map[string]*framework.FieldSchema{
					"path": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["audit_path"][0]),
					},
					"type": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["audit_type"][0]),
					},
					"description": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["audit_desc"][0]),
					},
					"options": &framework.FieldSchema{
						Type:        framework.TypeMap,
						Description: strings.TrimSpace(sysHelp["audit_opts"][0]),
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.WriteOperation:  b.handleEnableAudit,
					logical.DeleteOperation: b.handleDisableAudit,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["audit"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["audit"][1]),
			},

			&framework.Path{
				Pattern: "raw/(?P<path>.+)",

				Fields: map[string]*framework.FieldSchema{
					"path": &framework.FieldSchema{
						Type: framework.TypeString,
					},
					"value": &framework.FieldSchema{
						Type: framework.TypeString,
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation:   b.handleRawRead,
					logical.WriteOperation:  b.handleRawWrite,
					logical.DeleteOperation: b.handleRawDelete,
				},
			},
		},
	}
}

// SystemBackend implements logical.Backend and is used to interact with
// the core of the system. This backend is hardcoded to exist at the "sys"
// prefix. Conceptually it is similar to procfs on Linux.
type SystemBackend struct {
	Core *Core
}

// handleMountTable handles the "mounts" endpoint to provide the mount table
func (b *SystemBackend) handleMountTable(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Core.mounts.Lock()
	defer b.Core.mounts.Unlock()

	resp := &logical.Response{
		Data: make(map[string]interface{}),
	}
	for _, entry := range b.Core.mounts.Entries {
		info := map[string]string{
			"type":        entry.Type,
			"description": entry.Description,
		}
		resp.Data[entry.Path] = info
	}

	return resp, nil
}

// handleMount is used to mount a new path
func (b *SystemBackend) handleMount(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the options
	path := data.Get("path").(string)
	logicalType := data.Get("type").(string)
	description := data.Get("description").(string)

	if logicalType == "" {
		return logical.ErrorResponse(
				"backend type must be specified as a string"),
			logical.ErrInvalidRequest
	}

	// Create the mount entry
	me := &MountEntry{
		Path:        path,
		Type:        logicalType,
		Description: description,
	}

	// Attempt mount
	if err := b.Core.mount(me); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

// handleUnmount is used to unmount a path
//
// TODO: Run a rollback operation
// TODO: Return an error if things go bad.
// TODO: Think through error scenario: umount should clean up everything
//   and should fail if it can't. Perhaps add a "force" flag to YOLO it.
func (b *SystemBackend) handleUnmount(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	suffix := strings.TrimPrefix(req.Path, "mounts/")
	if len(suffix) == 0 {
		return logical.ErrorResponse("path cannot be blank"), logical.ErrInvalidRequest
	}

	// Attempt unmount
	if err := b.Core.unmount(suffix); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	return nil, nil
}

// handleRemount is used to remount a path
func (b *SystemBackend) handleRemount(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get the paths
	fromPath := data.Get("from").(string)
	toPath := data.Get("to").(string)
	if fromPath == "" || toPath == "" {
		return logical.ErrorResponse(
				"both 'from' and 'to' path must be specified as a string"),
			logical.ErrInvalidRequest
	}

	// Attempt remount
	if err := b.Core.remount(fromPath, toPath); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	return nil, nil
}

// handleRenew is used to renew a lease with a given VaultID
func (b *SystemBackend) handleRenew(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the options
	vaultID := data.Get("vault_id").(string)
	incrementRaw := data.Get("increment").(int)

	// Convert the increment
	increment := time.Duration(incrementRaw) * time.Second

	// Invoke the expiration manager directly
	resp, err := b.Core.expiration.Renew(vaultID, increment)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return resp, err
}

// handleRevoke is used to revoke a given VaultID
func (b *SystemBackend) handleRevoke(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the options
	vaultID := data.Get("vault_id").(string)

	// Invoke the expiration manager directly
	if err := b.Core.expiration.Revoke(vaultID); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

// handleRevokePrefix is used to revoke a prefix with many VaultIDs
func (b *SystemBackend) handleRevokePrefix(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the options
	prefix := data.Get("prefix").(string)

	// Invoke the expiration manager directly
	if err := b.Core.expiration.RevokePrefix(prefix); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

// handleAuthTable handles the "auth" endpoint to provide the auth table
func (b *SystemBackend) handleAuthTable(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Core.auth.Lock()
	defer b.Core.auth.Unlock()

	resp := &logical.Response{
		Data: make(map[string]interface{}),
	}
	for _, entry := range b.Core.auth.Entries {
		info := map[string]string{
			"type":        entry.Type,
			"description": entry.Description,
		}
		resp.Data[entry.Path] = info
	}
	return resp, nil
}

// handleEnableAuth is used to enable a new credential backend
func (b *SystemBackend) handleEnableAuth(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the options
	path := data.Get("path").(string)
	logicalType := data.Get("type").(string)
	description := data.Get("description").(string)

	if logicalType == "" {
		return logical.ErrorResponse(
				"backend type must be specified as a string"),
			logical.ErrInvalidRequest
	}

	// Create the mount entry
	me := &MountEntry{
		Path:        path,
		Type:        logicalType,
		Description: description,
	}

	// Attempt enabling
	if err := b.Core.enableCredential(me); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

// handleDisableAuth is used to disable a credential backend
func (b *SystemBackend) handleDisableAuth(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	suffix := strings.TrimPrefix(req.Path, "auth/")
	if len(suffix) == 0 {
		return logical.ErrorResponse("path cannot be blank"), logical.ErrInvalidRequest
	}

	// Attempt disable
	if err := b.Core.disableCredential(suffix); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

// handlePolicyList handles the "policy" endpoint to provide the enabled policies
func (b *SystemBackend) handlePolicyList(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the configured policies
	policies, err := b.Core.policy.ListPolicies()

	// Add the special "root" policy
	policies = append(policies, "root")
	return logical.ListResponse(policies), err
}

// handlePolicyRead handles the "policy/<name>" endpoint to read a policy
func (b *SystemBackend) handlePolicyRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)

	policy, err := b.Core.policy.GetPolicy(name)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	if policy == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"name":  name,
			"rules": policy.Raw,
		},
	}, nil
}

// handlePolicySet handles the "policy/<name>" endpoint to set a policy
func (b *SystemBackend) handlePolicySet(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	rules := data.Get("rules").(string)

	// Validate the rules parse
	parse, err := Parse(rules)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Override the name
	parse.Name = name

	// Update the policy
	if err := b.Core.policy.SetPolicy(parse); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

// handlePolicyDelete handles the "policy/<name>" endpoint to delete a policy
func (b *SystemBackend) handlePolicyDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	if err := b.Core.policy.DeletePolicy(name); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

// handleAuditTable handles the "audit" endpoint to provide the audit table
func (b *SystemBackend) handleAuditTable(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Core.audit.Lock()
	defer b.Core.audit.Unlock()

	resp := &logical.Response{
		Data: make(map[string]interface{}),
	}
	for _, entry := range b.Core.audit.Entries {
		info := map[string]interface{}{
			"type":        entry.Type,
			"description": entry.Description,
			"options":     entry.Options,
		}
		resp.Data[entry.Path] = info
	}
	return resp, nil
}

// handleEnableAudit is used to enable a new audit backend
func (b *SystemBackend) handleEnableAudit(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the options
	path := data.Get("path").(string)
	backendType := data.Get("type").(string)
	description := data.Get("description").(string)
	options := data.Get("options").(map[string]interface{})

	optionMap := make(map[string]string)
	for k, v := range options {
		vStr, ok := v.(string)
		if !ok {
			return logical.ErrorResponse("options must be string valued"),
				logical.ErrInvalidRequest
		}
		optionMap[k] = vStr
	}

	// Create the mount entry
	me := &MountEntry{
		Path:        path,
		Type:        backendType,
		Description: description,
		Options:     optionMap,
	}

	// Attempt enabling
	if err := b.Core.enableAudit(me); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

// handleDisableAudit is used to disable an audit backend
func (b *SystemBackend) handleDisableAudit(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)

	// Attempt disable
	if err := b.Core.disableAudit(path); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

// handleRawRead is used to read directly from the barrier
func (b *SystemBackend) handleRawRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	entry, err := b.Core.barrier.Get(path)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	if entry == nil {
		return nil, nil
	}
	resp := &logical.Response{
		Data: map[string]interface{}{
			"value": string(entry.Value),
		},
	}
	return resp, nil
}

// handleRawWrite is used to write directly to the barrier
func (b *SystemBackend) handleRawWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	value := data.Get("value").(string)
	entry := &Entry{
		Key:   path,
		Value: []byte(value),
	}
	if err := b.Core.barrier.Put(entry); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

// handleRawDelete is used to delete directly from the barrier
func (b *SystemBackend) handleRawDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	if err := b.Core.barrier.Delete(path); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

// sysHelp is all the help text for the sys backend.
var sysHelp = map[string][2]string{
	"mounts": {
		"List the currently mounted backends.",
		`
List the currently mounted backends: the mount path, the type of the backend,
and a user friendly description of the purpose for the mount.
		`,
	},

	"mount": {
		`Mount a new backend at a new path.`,
		`
Mount a backend at a new path. A backend can be mounted multiple times at
multiple paths in order to configure multiple separately configured backends.
Example: you might have an AWS backend for the east coast, and one for the
west coast.
		`,
	},

	"mount_path": {
		`The path to mount to. Example: "aws/east"`,
		"",
	},

	"mount_type": {
		`The type of the backend. Example: "passthrough"`,
		"",
	},

	"mount_desc": {
		`User-friendly description for this mount.`,
		"",
	},

	"remount": {
		"Move the mount point of an already-mounted backend.",
		`
Change the mount point of an already-mounted backend.
		`,
	},

	"remount_from": {
		"",
		"",
	},

	"remount_to": {
		"",
		"",
	},

	"renew": {
		"Renew a lease on a secret",
		`
When a secret is read, it may optionally include a lease interval
and a boolean indicating if renew is possible. For secrets that support
lease renewal, this endpoint is used to extend the validity of the
lease and to prevent an automatic revocation.
		`,
	},

	"vault_id": {
		"The vault identifier to renew. This is included with a lease.",
		"",
	},

	"increment": {
		"The desired increment in seconds to the lease",
		"",
	},

	"revoke": {
		"Revoke a leased secret immediately",
		`
When a secret is generated with a lease, it is automatically revoked
at the end of the lease period if not renewed. However, in some cases
you may want to force an immediate revocation. This endpoint can be
used to revoke the secret with the given Vault ID.
		`,
	},

	"revoke-prefix": {
		"Revoke all secrets generated in a given prefix",
		`
Revokes all the secrets generated under a given mount prefix. As
an example, "prod/aws/" might be the AWS logical backend, and due to
a change in the "ops" policy, we may want to invalidate all the secrets
generated. We can do a revoke prefix at "prod/aws/ops" to revoke all
the ops secrets. This does a prefix match on the Vault IDs and revokes
all matching leases.
		`,
	},

	"revoke-prefix-path": {
		`The path to revoke keys under. Example: "prod/aws/ops"`,
		"",
	},

	"auth-table": {
		"List the currently enabled credential backends.",
		`
List the currently enabled credential backends: the name, the type of the backend,
and a user friendly description of the purpose for the credential backend.
		`,
	},

	"auth": {
		`Enable a new credential backend with a name.`,
		`
Enable a credential mechanism at a new path. A backend can be mounted multiple times at
multiple paths in order to configure multiple separately configured backends.
Example: you might have an OAuth backend for GitHub, and one for Google Apps.
		`,
	},

	"auth_path": {
		`The path to mount to. Cannot be delimited. Example: "user"`,
		"",
	},

	"auth_type": {
		`The type of the backend. Example: "userpass"`,
		"",
	},

	"auth_desc": {
		`User-friendly description for this crential backend.`,
		"",
	},

	"policy-list": {
		`List the configured access control policies.`,
		`
List the names of the configured access control policies. Policies are associated
with client tokens to limit access to keys in the Vault.
		`,
	},

	"policy": {
		`Read, Modify, or Delete an access control policy.`,
		`
Read the rules of an existing policy, create or update the rules of a policy,
or delete a policy.
		`,
	},

	"policy-name": {
		`The name of the policy. Example: "ops"`,
		"",
	},

	"policy-rules": {
		`The rules of the policy. Either given in HCL or JSON format.`,
		"",
	},

	"audit-table": {
		"List the currently enabled audit backends.",
		`
List the currently enabled audit backends: the name, the type of the backend,
a user friendly description of the audit backend, and it's configuration options.
		`,
	},

	"audit_path": {
		`The name of the backend. Cannot be delimited. Example: "mysql"`,
		"",
	},

	"audit_type": {
		`The type of the backend. Example: "mysql"`,
		"",
	},

	"audit_desc": {
		`User-friendly description for this audit backend.`,
		"",
	},

	"audit_opts": {
		`Configuration options for the audit backend.`,
		"",
	},

	"audit": {
		`Enable or disable audit backends.`,
		`
Enable a new audit backend or disable an existing backend.
		`,
	},
}
