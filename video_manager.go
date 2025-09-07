package main

import (
	"log"
	"os/exec"
	"strings"
)

type VideoManager struct {
	Proxy    *ProxyManager
	VideoDir string
}

func NewVideoManager(proxy *ProxyManager, videoDir string) *VideoManager {
	return &VideoManager{Proxy: proxy, VideoDir: videoDir}
}

func (vm *VideoManager) Download(url string, descriptionUrls []string) (filename, description string, err error) {
	var printArg string
	useDescription := false
	for _, u := range descriptionUrls {
		if strings.Contains(url, u) {
			useDescription = true
			break
		}
	}
	if useDescription {
		printArg = "description"
	} else {
		printArg = "title"
	}

	var proxies []string
	if vm.Proxy != nil {
		if err := vm.Proxy.FetchProxies(); err == nil {
			proxies = vm.Proxy.Proxies
		}
	}
	var lastErr error
	tried := 0
	for _, proxy := range proxies {
		proxyArg := []string{"--proxy", proxy}
		ytArgs := append([]string{"-q", "--no-warnings", "--no-simulate", "-o", "%(id)s.%(ext)s", url, "-S", "ext", "--print", "filename", "--print", printArg}, proxyArg...)
		cmd := exec.Command("yt-dlp", ytArgs...)
		cmd.Dir = vm.VideoDir
		out, err := cmd.Output()
		tried++
		if err == nil {
			lines := strings.Split(string(out), "\n")
			filename = strings.TrimSpace(lines[0])
			if len(lines) > 1 {
				descRaw := strings.TrimSpace(strings.Join(lines[1:], "\n"))
				description = RemoveURLs(descRaw)
			}
			return filename, description, nil
		}
		lastErr = err
		log.Printf("yt-dlp failed with proxy %s: %v", proxy, err)
	}
	// Try without proxy if all proxies failed or no proxies
	ytArgs := []string{"-q", "--no-warnings", "--no-simulate", "-o", "%(id)s.%(ext)s", url, "-S", "ext", "--print", "filename", "--print", printArg}
	cmd := exec.Command("yt-dlp", ytArgs...)
	cmd.Dir = vm.VideoDir
	out, err := cmd.Output()
	if err != nil {
		log.Printf("yt-dlp command failed (no proxy): %v", err)
		log.Printf("yt-dlp command out: %v", string(out))
		if tried > 0 {
			return "", "", lastErr
		}
		return "", "", err
	}
	lines := strings.Split(string(out), "\n")
	filename = strings.TrimSpace(lines[0])
	if len(lines) > 1 {
		descRaw := strings.TrimSpace(strings.Join(lines[1:], "\n"))
		description = RemoveURLs(descRaw)
	}
	return filename, description, nil
}
