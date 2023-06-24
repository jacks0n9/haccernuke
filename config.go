package main

type NukeConfig struct {
	GuildID       string        `toml:"guildID"`
	Token         string        `toml:"token"`
	FeatureConfig FeatureConfig `toml:"feature_config"`
}
type FeatureConfig struct {
	AfterChannelConfig  AfterChannelConfig  `toml:"after_channels"`
	MemberRemovalConfig MemberRemovalConfig `toml:"member_removal"`
	DeleteRoles         bool                `toml:"delete_roles"`
	DeleteChannels      bool                `toml:"delete_channels"`
	AutoNuke            bool                `toml:"auto_nuke"`
	AutoAdmin           []string            `toml:"auto_admin"`
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
