//go:generate goversioninfo

package main

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/AllenDang/giu"
	"github.com/sqweek/dialog"
)

type launcherCommune struct {
	Patchline string `json:"last_patchline"`
	Username string `json:"last_username"`
	SelectedVersion int32 `json:"last_version"`
	LatestVersions map[string]int `json:"last_version_scan_result"`
	Mode string `json:"mode"`

	// authentication
	AuthTokens *accessTokens `json:"token"`
	Profiles *[]accountInfo `json:"profiles"`
	SelectedProfile int32 `json:"selected_profile"`

	// settings
	GameFolder string `json:"install_directory"`
	UserDataFolder string `json:"userdata_directory"`
	JreFolder string `json:"jre_directory"`
	UUID string `json:"uuid_override"`
}


const DEFAULT_USERNAME = "TransRights";
const DEFAULT_PATCHLINE = "release";

var (
	wMainWin *giu.MasterWindow
	wCommune = launcherCommune {
		Patchline: DEFAULT_PATCHLINE,
		Username: DEFAULT_USERNAME,
		LatestVersions: map[string]int{
			"release": 5,
			"pre-release": 12,
		},
		SelectedVersion: 4,
		Mode: "fakeonline",
		AuthTokens: nil,
		Profiles: nil,
		SelectedProfile: 0,

		GameFolder: DefaultGameFolder(),
		UserDataFolder: DefaultUserDataFolder(),
		JreFolder: DefaultJreFolder(),
		UUID: "",
	};
	wProgress = 0
	wDisabled = false
	wSelectedTab = 0
	wSelectedMode int32 = 0
	wSelectedPatchline int32 = 0

	wAdvanced = false
	w = false
)




func doAuthentication() {
	aTokens, err := getAuthTokens(wCommune.AuthTokens);

	if err != nil {
		showErrorDialog(fmt.Sprintf("Failed to get auth tokens: %s", err), "Auth failed.");
		wCommune.AuthTokens = nil;
		wCommune.Mode = "fakeonline";
		writeSettings();
		//loop.Do(updateWindow);
	}

	wCommune.AuthTokens = &aTokens;

	// get profile list ..
	authenticatedCheckForUpdatesAndGetProfileList();

}


func checkForUpdates() {
	if wCommune.Mode != "authenticated" {
		lastRelease := wCommune.LatestVersions["release"]
		lastPreRelease := wCommune.LatestVersions["pre-release"]

		latestRelease := findLatestVersionNoAuth(lastRelease, runtime.GOARCH, runtime.GOOS, "release");
		latestPreRelease := findLatestVersionNoAuth(lastPreRelease, runtime.GOARCH, runtime.GOOS, "pre-release");

		fmt.Printf("latestRelease: %d\n", latestRelease);
		fmt.Printf("latestPreRelease: %d\n", latestPreRelease);

		if latestRelease > lastRelease {
			fmt.Printf("Found new release version: %d\n", latestRelease);
			wCommune.LatestVersions["release"] = latestRelease;
		}

		if latestPreRelease > lastPreRelease {
			fmt.Printf("Found new pre-release version: %d\n", latestPreRelease);
			wCommune.LatestVersions["pre-release"] = latestPreRelease;
		}

		//loop.Run(updateWindow);
		writeSettings();
	}
}

func authenticatedCheckForUpdatesAndGetProfileList() {
	if wCommune.AuthTokens == nil {
		return;
	}
	if(wCommune.Mode != "authenticated") {
		return;
	}

	lData, err := getLauncherData(*wCommune.AuthTokens, runtime.GOARCH, runtime.GOOS);

	if err != nil {
		showErrorDialog(fmt.Sprintf("Failed to get launcher data: %s", err), "Auth failed.");
		wCommune.AuthTokens = nil;
		wCommune.Mode = "fakeonline";
		/*
		go func() {
			loop.Do(updateWindow);
		}(); */
		writeSettings();
	}

	lastReleaseVersion := wCommune.LatestVersions["release"];
	latestReleaseVersion := lData.Patchlines.Release.Newest;

	lastPreReleaseVersion := wCommune.LatestVersions["pre-release"];
	latestPreReleaseVersion := lData.Patchlines.PreRelease.Newest;

	if latestReleaseVersion > lastReleaseVersion {
		fmt.Printf("found new release: %d\n", lastReleaseVersion)
		wCommune.LatestVersions["release"] = latestReleaseVersion;
	}
	if latestPreReleaseVersion > lastPreReleaseVersion {
		fmt.Printf("found new release: %d\n", lastPreReleaseVersion)
		wCommune.LatestVersions["pre-release"] = latestPreReleaseVersion;
	}

	wCommune.Profiles = &lData.Profiles;

	writeSettings();
}

func reAuthenticate() {
	if wCommune.AuthTokens != nil && wCommune.Mode == "authenticated" {
		aTokens, err:= getAuthTokens(*wCommune.AuthTokens);

		if err != nil {
			showErrorDialog(fmt.Sprintf("Failed to authenticate: %s", err), "Auth failed.");
			wCommune.AuthTokens = nil;
			wCommune.Mode = "fakeonline";
			writeSettings();
		}

		wCommune.AuthTokens = &aTokens;
		authenticatedCheckForUpdatesAndGetProfileList();
	}
}

func writeSettings() {
	fmt.Printf("Saving settings ...\n");
	jlauncher, _ := json.Marshal(wCommune);

	err := os.MkdirAll(filepath.Dir(getLauncherJson()), 0666);
	if err != nil {
		fmt.Printf("error writing settings: %s\n", err);
		return;
	}

	err = os.WriteFile(getLauncherJson(), jlauncher, 0666);
	if err != nil {
		fmt.Printf("error writing settings: %s\n", err);
		return;
	}
}

func getDefaultSettings() {
	writeSettings();
	go checkForUpdates();

}

func getLauncherJson() string {
	return filepath.Join(LauncherFolder(), "launcher.json");
}

func readSettings() {
	_, err := os.Stat(getLauncherJson())
	if err != nil {
		getDefaultSettings();
	} else {
		data, err := os.ReadFile(getLauncherJson());
		if err != nil{
			getDefaultSettings();
			return;
		}
		json.Unmarshal(data, &wCommune);

		if wCommune.GameFolder != GameFolder() {
			wCommune.GameFolder = GameFolder();
		}

		fmt.Printf("Reading last settings: \n");
		fmt.Printf("username: %s\n", wCommune.Username);
		fmt.Printf("patchline: %s\n", wCommune.Patchline);
		fmt.Printf("last used version: %d\n", wCommune.SelectedVersion);
		fmt.Printf("newest known release: %d\n", wCommune.LatestVersions["release"])
		fmt.Printf("newest known pre-release: %d\n", wCommune.LatestVersions["pre-release"])

	}
}


func valToChannel(vchl int) string {
	switch vchl {
		case 0:
			return "release";
		case 1:
			return "pre-release";
		default:
			return "release";
	}
}

func channelToVal(channel string) int {
	switch channel {
		case "release":
			return 0;
		case "pre-release":
			return 1;
		default:
			return 0;
	}
}

func startGame() {
	// disable the current window
	wDisabled = true;

	// enable the window again once done
	defer func() {
		wDisabled = false;
	}();

	err := installJre(updateProgress);

	if err != nil {
		showErrorDialog(fmt.Sprintf("Error getting the JRE: %s", err), "Install JRE failed.");
		return;
	};

	err = installGame(int(wCommune.SelectedVersion), wCommune.Patchline, updateProgress);

	if err != nil {
		showErrorDialog(fmt.Sprintf("Error getting the game: %s", err), "Install game failed.");
		return;
	};

	err = launchGame(int(wCommune.SelectedVersion), wCommune.Patchline, wCommune.Username, getUUID());

	if err != nil {
		showErrorDialog(fmt.Sprintf("Error running the game: %s", err), "Run game failed.");
		return;
	};
}

func patchLineMenu() giu.Widget {

	return giu.Layout{
		giu.Label("Patchline: "),
		giu.Combo("#patchline", wCommune.Patchline, []string{"release", "pre-release"}, &wSelectedPatchline).OnChange(func() {
			wCommune.Patchline = valToChannel(int(wSelectedPatchline));
			wCommune.SelectedVersion = int32(wCommune.LatestVersions[wCommune.Patchline]);
		}),
	}
}


func versionMenu() giu.Widget {
	versions := []string {};

	latest := wCommune.LatestVersions[wCommune.Patchline];

	for i := range latest {
		txt := "Version "+strconv.Itoa(i+1);
		if isGameVersionInstalled(i+1, wCommune.Patchline) {
			txt += " - installed";
		} else {
			txt += " - not installed";
		}

		versions = append(versions, txt);
	}

	gotVersion := int(wCommune.SelectedVersion);
	selectedChannel := wCommune.Patchline;

	disabled := !isGameVersionInstalled(gotVersion, selectedChannel) || wDisabled;


	return giu.Layout{
		giu.Label("Version: "),

		giu.Row(
			giu.Combo("#version", versions[gotVersion-1], versions, &wCommune.SelectedVersion).Flags(giu.ComboFlagsNone),
			giu.Button("Delete").Disabled(disabled).OnClick(func() {
				wDisabled = true;

				go func() {
						installDir := getVersionInstallPath(gotVersion, wCommune.Patchline);
						err := os.RemoveAll(installDir);
						if err != nil {
							showErrorDialog(fmt.Sprintf("failed to remove: %s", err), "failed to remove");
						}

						wDisabled = false;
					}();
				},
			),
		),
	};
}


func labeledTextInput(label string, value *string, disabled bool) giu.Widget {
	if value == nil {
		panic("failed to initalize browse button");
	}

	flags := giu.InputTextFlagsNone;

	if wDisabled || disabled {
		flags = giu.InputTextFlagsReadOnly;
	}

	return giu.Layout{
		giu.Label(label+": "),
		giu.InputText(value).Flags(flags).Hint(label),
	}
}


func drawDivider(label string) giu.Widget {
	return giu.Row(
		giu.Label(label),
		       giu.Separator(),
	);
}


func browseButton(label string, value *string) giu.Widget {
	if value == nil {
		panic("failed to initalize browse button");
	}

	return giu.Layout{
		giu.Label(label + ": "),
		giu.Row(
			giu.InputText(value).Hint(label).Label(label + ":  "),
			giu.Button("Browse").OnClick(func() {
				dir, err := dialog.Directory().Title("Select "+label).Browse();
				if err != nil {
					if err != dialog.ErrCancelled {
						showErrorDialog(fmt.Sprintf("Failed: %s", err), "Error reading directory");
					}
				}

				*value = dir;
			}),
		),
	}
}



func modeSelector () giu.Widget {
	return giu.Layout{
		giu.Label("Launch Mode: "),
		giu.Combo("#launchMode", wCommune.Mode, []string {"Offline Mode", "Fake Online Mode", "Authenticated"}, &wSelectedMode).OnChange(func() {
			switch(wSelectedMode) {
				case 0:
					wCommune.Mode = "offline";
				case 1:
					wCommune.Mode = "fakeonline";
				case 2:
					wCommune.Mode = "authenticated";
			}
		}),
	};
}



func drawAuthenticatedSettings() giu.Widget {

	if wCommune.Mode != "authenticated" {
		return giu.Custom(func() {});
	}

	logoutDisabled := wDisabled || (wCommune.AuthTokens == nil);
	loginDisabled := wDisabled || (wCommune.AuthTokens != nil);
	profileList := []string{};

	if wCommune.Profiles != nil {
		for _, profile := range *wCommune.Profiles {
			profileList = append(profileList, profile.Username);
		}
	}
	//profilesDisabled := wDisabled || wCommune.Profiles == nil;

	return giu.Layout{
		drawDivider("Authentication"),
		giu.Row(
			giu.Button("Login (OAuth 2.0").Disabled(loginDisabled).OnClick(func() {
				go doAuthentication();
			}),
			giu.Button("Logout").Disabled(logoutDisabled).OnClick(func() {
				wCommune.AuthTokens = nil;
				wCommune.Profiles = nil;
				writeSettings();
			}),
			giu.Label("Select profile"),
			giu.Combo("##selectProfile", profileList[wCommune.SelectedProfile], profileList, &wCommune.SelectedProfile),
		),
	};

}


func updateProgress(done int64, total int64) {
	lastProgress := wProgress;
	newProgress := int((float64(done) / float64(total)) * 100.0);

	if newProgress != lastProgress {
		wProgress = newProgress;
	}
}

func createDownloadProgress () giu.Widget {
	return giu.ProgressBar(float32(wProgress));
}

func drawStartGame() giu.Widget{
	return &giu.Layout {
			labeledTextInput("Username", &wCommune.Username, wCommune.Mode == "authenticated"),
			modeSelector(),
			drawDivider("Version"),
			patchLineMenu(),
			versionMenu(),
			createDownloadProgress(),
			giu.Button("Start Game").Disabled(wDisabled).OnClick(func() {
				go startGame();
			}),
			//drawAuthenticatedSettings(),
	}
}

func drawSettings() giu.Widget{
	return giu.Layout{
		drawDivider("Directories"),
		browseButton("Game Location", &wCommune.GameFolder),
		browseButton("JRE Location", &wCommune.JreFolder),
		browseButton("UserData Location", &wCommune.UserDataFolder),
		drawDivider("Advanced"),
		labeledTextInput("★UUID Override", &wCommune.UUID, wCommune.Mode == "authenticated"),
	};

	/*
	return &goey.VBox{
		AlignMain: goey.MainStart,
		Children: []base.Widget {
			drawDivider("Directories"),
			browseButton("Game Location", &wCommune.GameFolder),
			browseButton("JRE Location", &wCommune.JreFolder),
			browseButton("Game UserData Location", &wCommune.UserDataFolder),
			drawDivider("Advanced"),
			labeledTextInput("★UUID Override", &wCommune.UUID, wCommune.Mode == "authenticated"),
		},
	};*/
}


func drawWidgets() {

	giu.SingleWindow().Layout(
		giu.TabBar().TabItems(
			giu.TabItem("Game").Layout(
				drawStartGame(),
			),
		),
		giu.TabBar().TabItems(
			giu.TabItem("Settings").Layout(
				drawSettings(),
			),
		),

	)
	/*
	return &goey.Tabs {
		Value: wSelectedTab,
		OnChange: func( v int ) {wSelectedTab = v;},
		Children: []goey.TabItem {
			{
				Caption: "Game",
				Child: drawStartGame(),
			},
			{
				Caption: "Settings",
				Child: drawSettings(),
			},
		},
	};*/

}

func createWindow() error {

	win := giu.NewMasterWindow("HytaleSP", 400, 200, giu.MasterWindowFlagsNotResizable)
	if win != nil {
		return fmt.Errorf("result from NewMasterWindow was nil");
	}

	win.SetCloseCallback(func() bool {
		writeSettings();
		return false;
	});


	f, err := embeddedImages.Open(path.Join("Resources", "icon.png"));
	if err != nil {
		return nil;
	}
	defer f.Close()

	image, _, err := image.Decode(f)
	win.SetIcon(image);

	// create the ui
	win.Run(drawWidgets);

	wMainWin = win;

	return nil
}


func showErrorDialog(msg string, title string) {
		dlg := dialog.Message(msg);
		dlg.Title(title);
		dlg.Error();
}




func main() {

	os.MkdirAll(MainFolder(), 0775);
	os.MkdirAll(LauncherFolder(), 0775);
	os.MkdirAll(ServerDataFolder(), 0775);
	readSettings();

	os.MkdirAll(UserDataFolder(), 0775);
	os.MkdirAll(JreFolder(), 0775);
	os.MkdirAll(GameFolder(), 0775);

	go reAuthenticate();
	go checkForUpdates();

	err := createWindow();
	if err != nil {
		showErrorDialog(fmt.Sprintf("Error occured while creating window: %s", err), "Error while creating window");
	}

}
