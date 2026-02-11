package api

// Device represents a LaMetric device info response.
type Device struct {
	ID           string    `json:"id"`
	UUID         string    `json:"uuid"`
	Name         string    `json:"name"`
	SerialNumber string    `json:"serial_number"`
	OSVersion    string    `json:"os_version"`
	Mode         string    `json:"mode"`
	Model        string    `json:"model"`
	Audio        Audio     `json:"audio"`
	Bluetooth    Bluetooth `json:"bluetooth"`
	Display      Display   `json:"display"`
	WiFi         WiFi      `json:"wifi"`
}

// Display represents display settings.
type Display struct {
	Brightness     int          `json:"brightness"`
	BrightnessMode string       `json:"brightness_mode"` // auto|manual
	Width          int          `json:"width"`
	Height         int          `json:"height"`
	Type           string       `json:"type"`
	Screensaver    *Screensaver `json:"screensaver,omitempty"`
}

// Screensaver represents screensaver configuration.
type Screensaver struct {
	Enabled bool   `json:"enabled"`
	Widget  string `json:"widget,omitempty"`
}

// Audio represents audio settings.
type Audio struct {
	Volume int `json:"volume"`
}

// Bluetooth represents bluetooth settings.
type Bluetooth struct {
	Available    bool   `json:"available"`
	Active       bool   `json:"active"`
	Name         string `json:"name"`
	MAC          string `json:"mac,omitempty"`
	Pairable     bool   `json:"pairable"`
	Discoverable bool   `json:"discoverable"`
}

// WiFi represents wifi status (read-only).
type WiFi struct {
	Active     bool   `json:"active"`
	MAC        string `json:"mac"`
	SSID       string `json:"ssid"`
	IP         string `json:"ip"`
	Netmask    string `json:"netmask"`
	Strength   int    `json:"strength"`
	Encryption string `json:"encryption,omitempty"`
}

// NotificationRequest is the payload sent to create a notification.
type NotificationRequest struct {
	Priority string            `json:"priority,omitempty"`  // info|warning|critical
	IconType string            `json:"icon_type,omitempty"` // none|info|alert
	Lifetime int               `json:"lifetime,omitempty"`  // ms
	Model    NotificationModel `json:"model"`
}

// NotificationModel holds the frames, sound, and cycle config.
type NotificationModel struct {
	Frames []Frame `json:"frames"`
	Sound  *Sound  `json:"sound,omitempty"`
	Cycles int     `json:"cycles,omitempty"`
}

// Frame represents a single notification frame.
type Frame struct {
	Icon      string    `json:"icon,omitempty"`
	Text      string    `json:"text,omitempty"`
	GoalData  *GoalData `json:"goalData,omitempty"`
	ChartData []int     `json:"chartData,omitempty"`
}

// GoalData represents goal progress data.
type GoalData struct {
	Start   int    `json:"start"`
	Current int    `json:"current"`
	End     int    `json:"end"`
	Unit    string `json:"unit,omitempty"`
}

// Sound represents a notification sound.
type Sound struct {
	Category string `json:"category,omitempty"` // notifications|alarms
	ID       string `json:"id,omitempty"`
	Repeat   int    `json:"repeat,omitempty"`
}

// Notification represents a created/queued notification.
type Notification struct {
	ID       string `json:"id"`
	Type     string `json:"type,omitempty"`
	Priority string `json:"priority,omitempty"`
	Created  string `json:"created,omitempty"`
}

// App represents an installed app/widget.
type App struct {
	Package string            `json:"package"`
	Vendor  string            `json:"vendor,omitempty"`
	Version string            `json:"version,omitempty"`
	Widgets map[string]Widget `json:"widgets,omitempty"`
}

// Widget represents a single widget within an app.
type Widget struct {
	ID       string `json:"id,omitempty"`
	Index    int    `json:"index"`
	Package  string `json:"package,omitempty"`
	Settings any    `json:"settings,omitempty"`
}

// DisplayUpdate is the payload for updating display settings.
type DisplayUpdate struct {
	Brightness     *int    `json:"brightness,omitempty"`
	BrightnessMode *string `json:"brightness_mode,omitempty"`
}

// AudioUpdate is the payload for updating audio settings.
type AudioUpdate struct {
	Volume *int `json:"volume,omitempty"`
}

// BluetoothUpdate is the payload for updating bluetooth settings.
type BluetoothUpdate struct {
	Active       *bool   `json:"active,omitempty"`
	Name         *string `json:"name,omitempty"`
	Pairable     *bool   `json:"pairable,omitempty"`
	Discoverable *bool   `json:"discoverable,omitempty"`
}
