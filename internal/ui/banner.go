package ui

import "github.com/charmbracelet/lipgloss"

const bannerRaw = `
 ██████╗ ██╗   ██╗██╗     ███████╗███████╗
 ██╔══██╗██║   ██║██║     ██╔════╝██╔════╝
 ██████╔╝██║   ██║██║     ███████╗█████╗
 ██╔═══╝ ██║   ██║██║     ╚════██║██╔══╝
 ██║     ╚██████╔╝███████╗███████║███████╗
 ╚═╝      ╚═════╝ ╚══════╝╚══════╝╚══════╝`

var bannerStyle = lipgloss.NewStyle().
	Foreground(Purple).
	Bold(true)

var taglineStyle = lipgloss.NewStyle().
	Foreground(Cyan).
	Italic(true).
	MarginTop(1).
	MarginBottom(1)

// Banner returns the styled Pulse banner.
func Banner() string {
	return bannerStyle.Render(bannerRaw) + "\n" +
		taglineStyle.Render("  P2P file sharing. No servers. No cloud. Just peers.") + "\n"
}

// Version returns the formatted version string.
func Version() string {
	return Muted.Render("Pulse v2.0")
}
