export interface RegisterRequestPayload {
  username: string;
  email: string;
  master_key_hash: string;
  master_salt: string;
  code: string;
}

export interface LoginRequestPayload {
  identifier: string;
  master_key_hash: string;
}

export interface VaultItem {
	ID: string;
	UserID: string;
	EncryptedData: string;
	Category: string;
	CreatedAt: string;
	UpdatedAt: string;
}

export interface DecryptedVaultItem {
	ID: string;
	UserID: string;
	Category: string;
	CreatedAt: string;
	UpdatedAt: string;
	name: string;
	account: string;
	website?: string;
	password?: string;
	notes?: string;
}