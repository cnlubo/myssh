package myssh

import (
	"fmt"
	"github.com/cnlubo/myssh/confirmation"
	"github.com/cnlubo/myssh/prompt"
	"github.com/cnlubo/myssh/selection"
	"github.com/cnlubo/myssh/utils"
	"github.com/cnlubo/promptx"
	"github.com/fatih/color"
	"github.com/muesli/termenv"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	// DefaultKey is the default alias name of SSH key
	DefaultKey = "default"
	// HookName is the name of a hook that is called when present after using a key
	HookName       = "hook"
	DefaultBitSize = 2048
	customTemplate = `
{{- if .Prompt -}}
  {{ Bold .Prompt }}
{{ end -}}
{{ if .IsFiltered }}
  {{- print .FilterPrompt " " .FilterInput }}
{{ end }}

{{- range  $i, $choice := .Choices }}
  {{- if IsScrollUpHintPosition $i }}
    {{- print "⇡ " -}}
  {{- else if IsScrollDownHintPosition $i -}}
    {{- print "⇣ " -}} 
  {{- else -}}
    {{- print "  " -}}
  {{- end -}} 

  {{- if eq $.SelectedIndex $i }}
   {{- print "[" (Foreground "32" (Bold "x")) "] " (Selected $choice) "\n" }}
  {{- else }}
    {{- print "[ ] " (Unselected $choice) "\n" }}
  {{- end }}
{{- end}}`
	resultTemplate = `
		{{- print .Prompt " " (Foreground "32"  (name .FinalChoice)) "\n" -}}
		`
)

// SSHKey struct includes both private/public keys & isDefault flag
type SSHKey struct {
	PublicKey  string
	PrivateKey string
	IsDefault  bool
	Type       *KeyType
	KeyDesc    string
	Alias      string
}

func KeyStoreInit(env *Environment) error {

	keyStorePath := env.SKMPath

	// create keystore dir
	if found := utils.PathExist(keyStorePath); !found {
		err := os.Mkdir(keyStorePath, 0755)
		if err != nil {
			return errors.Wrap(err, "Create keyStore dir failed")
		}
	} else {
		if ok, _ := utils.IsEmpty(keyStorePath); !ok {
			return errors.New(fmt.Sprintf("KeyStore: %s is not empty", keyStorePath))
		}
	}

	if err := createDefaultSSHKey(env); err != nil {
		return err
	}
	fmt.Println()
	utils.PrintN(utils.Info, "SSH key store initialized!\n")
	utils.PrintN(utils.Info, fmt.Sprintf("Key store location is:%s\n", keyStorePath))
	return nil
}

func createDefaultSSHKey(env *Environment) error {
	var (
		keyFile           string
		defaultPrivateKey string
		defaultPublicKey  string
		confirm           bool
		keyStorePath      = env.SKMPath
		sshKeyPath        = env.SSHPath
		foundDefaultKey   = false
		err               error
	)

	// check existing keys in SSHPath (~/.ssh)
	for _, kt := range SupportedKeyTypes {
		keyFile = filepath.Join(sshKeyPath, kt.PrivateKey())
		if exist := utils.PathExist(keyFile); exist {
			foundDefaultKey = true
			defaultPrivateKey = kt.PrivateKey()
			defaultPublicKey = kt.PublicKey()
			break
		}
	}

	if foundDefaultKey {
		// Create default alias directory
		err = os.Mkdir(filepath.Join(keyStorePath, DefaultKey), 0755)
		if err != nil {
			return errors.Wrap(err, "create default alias directory failed")
		}
		// Move ~/.ssh/ key to default keystore dir
		err = os.Rename(keyFile, filepath.Join(keyStorePath, DefaultKey, defaultPrivateKey))
		if err != nil {
			return errors.Wrap(err, "move default PrivateKey file failed")
		}
		err = os.Rename(filepath.Join(sshKeyPath, defaultPublicKey), filepath.Join(keyStorePath, DefaultKey, defaultPublicKey))
		if err != nil {
			return errors.Wrap(err, "move default PublicKey file failed")
		}
	} else {
		input := confirmation.New("Do you want to create default SSHKey",
			confirmation.NewValue(false))
		//input.Template = confirmation.TemplateYN
		//input.ResultTemplate = confirmation.ResultTemplateYN
		//input.KeyMap.SelectYes = append(input.KeyMap.SelectYes, "+")
		//input.KeyMap.SelectNo = append(input.KeyMap.SelectNo, "-")
		input.ResultTemplate = ""
		confirm, err = input.RunPrompt()
		utils.CheckAndExit(err)
		if confirm {
			//// Create default alias directory
			//err := os.Mkdir(filepath.Join(keyStorePath, DefaultKey), 0755)
			//if err != nil {
			//	return errors.Wrap(err, "create default alias directory failed")
			//}
			//// create default SSHKey
			//ag := &createOptions{
			//	Alias:   DefaultKey,
			//	Comment: DefaultKey,
			//}
			//
			// supported keyType
			var sshKeyType []KeyType
			for _, kt := range SupportedKeyTypes {
				sshKeyType = append(sshKeyType, kt)
			}
			//
			//cfg := &promptx.SelectConfig{
			//	ActiveTpl:    `»  {{ .Name | cyan }}`,
			//	InactiveTpl:  `  {{ .Name | white }}`,
			//	SelectPrompt: "SSH Key Type",
			//	SelectedTpl:  `{{ "» " | green }}{{ "KeyType:" | green }}{{ .Name | green }}`,
			//	DisPlaySize:  9,
			//	DetailsTpl: `
			//--------- SSH Key Type ----------
			//{{ "Name:" | faint }} {{ .Name | faint }}
			//{{ "KeyBaseName:" | faint }} {{ .KeyBaseName | faint }}
			//{{ "SupportsVariableBitsize:" | faint }} {{ .SupportsVariableBitsize }}`,
			//}
			//
			//s := &promptx.Select{
			//	Items:  sshKeyType,
			//	Config: cfg,
			//}
			//idx := s.Run()
			//ag.Type = sshKeyType[idx].Name
			//if sshKeyType[idx].SupportsVariableBitsize {
			//	ag.Bits = DefaultBitSize
			//}
			//_, err = createKey(env, ag)
			//if err != nil {
			//	return err
			//}
			//// Set key with related alias as default used key
			//keyMap, _ := LoadSSHKeys(env)
			//err = createLink(ag.Alias, keyMap, env)
			//if err != nil {
			//	return err
			//}
			//// Run a potential hook
			//runHook(ag.Alias, env)
			//utils.PrintN(utils.Info, fmt.Sprintf("Now using SSH key: [%s]", ag.Alias))

			//sp := selection.New("Select SSH Key Type:",
			//	selection.Choices([]string{"Horse", "Car", "Plane", "Bike"}))

			type article struct {
				ID   string
				Name string
			}

			//choices := []article{
			//	{ID: "123", Name: "Article A"},
			//	{ID: "321", Name: "Article B"},
			//	{ID: "345", Name: "Article C"},
			//	{ID: "456", Name: "Article D"},
			//	{ID: "444", Name: "Article E"},
			//}
			//blue := termenv.String().Foreground(termenv.ANSI256Color(32)) // nolint:gomnd
			//blue := color.New(color.BgHiBlue, color.Bold).PrintfFunc()
			blue := color.New(color.Bold, color.FgHiBlue).SprintfFunc()
			sp := selection.New("Select SSH Key Type:",
				selection.Choices(sshKeyType))
			sp.PageSize = 5
			sp.Filter = nil
			sp.Template = customTemplate
			sp.ResultTemplate = resultTemplate
			sp.SelectedChoiceStyle = func(c *selection.Choice) string {
				a, _ := c.Value.(KeyType)

				//return blue.Bold().Styled(a.Name) + " " + termenv.String("("+a.KeyBaseName+")").Faint().String()
				return blue(a.Name)
			}
			sp.UnselectedChoiceStyle = func(c *selection.Choice) string {
				a, _ := c.Value.(KeyType)

				return a.Name + " " + termenv.String("("+a.KeyBaseName+")").Faint().String()
			}
			sp.ExtendedTemplateFuncs = map[string]interface{}{
				"name": func(c *selection.Choice) string { return c.Value.(KeyType).Name },
			}

			choice, err := sp.RunPrompt()
			if err != nil {
				fmt.Printf("Error: %v\n", err)

				os.Exit(1)
			}
			fmt.Println(choice)

		} else {
			return errors.New("Exit...")
		}
	}
	return nil
}

type createOptions struct {
	Alias   string
	Bits    int
	Comment string
	Type    string
}

// bits 密钥长度（bits），RSA最小要求768位，默认是2048位；DSA密钥必须是1024位(FIPS 1862标准规定)
// comment 提供一个注释
// type 指定生成密钥的类型 当前支持rsa 和ed25519

func CreateSSHKey(env *Environment) (map[string]*SSHKey, error) {

	keyMap, err := LoadSSHKeys(env)
	// alias name
	p := promptx.NewDefaultPrompt(func(line []rune) error {
		if strings.TrimSpace(string(line)) == "" {
			return inputEmptyErr
		} else if len(line) > 12 {
			return inputTooLongErr
		}
		if len(keyMap) > 0 {
			if _, ok := keyMap[string(line)]; ok {
				return errors.New(fmt.Sprintf("alias [%s] already exists", string(line)))
			}
		}
		return nil

	}, "AliasName:")

	AliasName, err := p.Run()
	if err != nil {
		return nil, err
	}

	ag := &createOptions{
		Alias: AliasName,
	}

	// supported keyType
	var sshKeyType []KeyType
	for _, kt := range SupportedKeyTypes {
		sshKeyType = append(sshKeyType, kt)

	}

	cfg := &promptx.SelectConfig{
		ActiveTpl:    `»  {{ .Name | cyan }}`,
		InactiveTpl:  `  {{ .Name | white }}`,
		SelectPrompt: "SSH Key Type",
		SelectedTpl:  `{{ "» " | green }}{{ "KeyType:" | green }}{{ .Name | green }}`,
		DisPlaySize:  9,
		DetailsTpl: `
--------- SSH Key Type ----------
{{ "Name:" | faint }} {{ .Name | faint }}
{{ "KeyBaseName:" | faint }} {{ .KeyBaseName | faint }}
{{ "SupportsVariableBitsize:" | faint }} {{ .SupportsVariableBitsize }}`,
	}

	s := &promptx.Select{
		Items:  sshKeyType,
		Config: cfg,
	}
	idx := s.Run()
	ag.Type = sshKeyType[idx].Name

	// BitSize
	if sshKeyType[idx].SupportsVariableBitsize {
		var bitSize int
		p = promptx.NewDefaultPrompt(func(line []rune) error {
			if strings.TrimSpace(string(line)) != "" {
				if _, err := strconv.Atoi(string(line)); err != nil {
					return notNumberErr
				}
			}
			return nil

		}, "BitSize(default:2048):")

		p.Default = strconv.Itoa(DefaultBitSize)
		bit, err := p.Run()
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(bit) == "" {
			bitSize = DefaultBitSize
		} else {
			bitSize, _ = strconv.Atoi(bit)
		}
		ag.Bits = bitSize
	}
	// Comment
	p = promptx.NewDefaultPrompt(func(line []rune) error {
		// allow empty
		return nil

	}, "Comment:")

	comment, err := p.Run()
	if err != nil {
		return nil, err
	}
	ag.Comment = comment

	if sshkeys, err := createKey(env, ag); err != nil {
		return nil, err
	} else {
		return sshkeys, nil
	}
}

func createKey(env *Environment, options *createOptions) (map[string]*SSHKey, error) {

	// check aliasName
	keyMap, err := LoadSSHKeys(env)
	if len(keyMap) > 0 {
		if _, ok := keyMap[options.Alias]; ok {
			return nil, errors.New(fmt.Sprintf("alias [%s] already exists", options.Alias))
		}
	}

	// Create alias directory
	aliasPath := filepath.Join(env.SKMPath, options.Alias)
	if exists := utils.PathExist(aliasPath); exists {
		err = os.RemoveAll(aliasPath)
		if err != nil {
			return nil, errors.Wrap(err, "remove alias dir failed")
		}
	}
	err = os.Mkdir(aliasPath, 0755)
	if err != nil {
		return nil, errors.Wrap(err, "create alias dir failed")
	}

	// generate args
	var args []string
	args = append(args, "-t")
	args = append(args, options.Type)

	keyTypeSettings, ok := SupportedKeyTypes[options.Type]
	if !ok {
		return nil, errors.Wrapf(err, "%s is not a supported KeyType.", options.Type)
	}

	args = append(args, "-f")
	fileName := keyTypeSettings.KeyBaseName
	args = append(args, filepath.Join(env.SKMPath, options.Alias, fileName))

	if keyTypeSettings.SupportsVariableBitsize {
		args = append(args, "-b")
		if options.Bits == 0 {
			args = append(args, strconv.Itoa(DefaultBitSize))
		} else {
			args = append(args, strconv.Itoa(options.Bits))
		}
	}

	if options.Comment != "" {
		args = append(args, "-C")
		args = append(args, options.Comment)
	}

	if ok := utils.Execute("", "ssh-keygen", args...); !ok {
		return nil, errors.New(fmt.Sprintf("SSHkey [%s] created failed\n", options.Alias))
	}
	utils.PrintN(utils.Info, fmt.Sprintf("SSHkey [%s] created!\n", options.Alias))

	sshKey := loadSingleKey(filepath.Join(env.SKMPath, options.Alias), options.Alias, env)
	keys := map[string]*SSHKey{}
	if sshKey != nil {
		keys[options.Alias] = sshKey
	}
	return keys, nil
}

func CopySSHKey(connectStr string, port string, env *Environment) error {
	var (
		identityfile string
		keyList      []*SSHKey
	)

	keyMap, _ := LoadSSHKeys(env)
	fmt.Println()
	cfg := &promptx.SelectConfig{
		ActiveTpl:    `» {{ .Alias | cyan }} {{"[" | cyan }}{{.Type.Name | cyan }}{{"]" | cyan }}`,
		InactiveTpl:  `  {{ .Alias | white }} {{"[" | white }}{{.Type.Name | white }}{{"]" | white }}`,
		SelectPrompt: "one ssh key",
		SelectedTpl:  `» {{ .Alias | green | bold }}`,
		DisPlaySize:  9,
		DetailsTpl: `
--------- SSH Key ----------
{{ "Alias:" | blue | faint }} {{ .Alias | blue | faint }}
{{ "Type:" | blue | faint }} {{ .Type.Name | blue | faint }}`,
	}

	for k := range keyMap {
		keyList = append(keyList, keyMap[k])
	}
	s := &promptx.Select{
		Items:  keyList,
		Config: cfg,
	}
	idx := s.Run()
	key, _ := keyMap[keyList[idx].Alias]
	alias := keyList[idx].Alias
	identityfile = filepath.Join(env.SKMPath, alias, key.Type.PrivateKey())
	// check privateKey
	_, err := privateKeyFile(identityfile, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("bad identityfile: [%s]", identityfile))
	}
	// copy SSHKey
	result := CopyKey(connectStr, port, identityfile)
	if !result {
		return errors.New("Copy SSHKey failed")
	}
	return nil
}

func CopyKey(connectStr string, port string, keyPath string) bool {

	var args []string
	args = append(args, "-p")
	args = append(args, port)
	args = append(args, "-i")
	args = append(args, keyPath)
	args = append(args, connectStr)

	return utils.Execute("", "ssh-copy-id", args...)
}

// loads all the SSH keys from key store
func LoadSSHKeys(env *Environment) (map[string]*SSHKey, error) {

	keys := map[string]*SSHKey{}

	// Walkthrough SSH key store and load all the keys
	err := filepath.Walk(env.SKMPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		if path == env.SKMPath {
			return nil
		}

		if f.IsDir() {
			// Load private/public keys
			key := loadSingleKey(path, f.Name(), env)

			if key != nil {
				keys[f.Name()] = key
			}
		}

		return nil
	})

	if err != nil {
		return keys, fmt.Errorf("filepath.Walk() returned %v\n", err)
	}

	return keys, nil
}

func loadSingleKey(keyPath string, alias string, env *Environment) *SSHKey {

	key := &SSHKey{}
	// Walkthrough SSH key store and load all the keys
	err := filepath.Walk(keyPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		if path == keyPath {
			return nil
		}

		if f.IsDir() {
			return nil
		}

		if strings.Contains(f.Name(), ".pub") {
			key.PublicKey = path
			return nil
		}

		kt, ok := SupportedKeyTypes.GetByFilename(f.Name())
		if !ok {
			return nil
		}
		key.Type = &kt
		key.Alias = alias

		// Check if key is in use
		key.PrivateKey = path

		if path == utils.ParseOriginalFilePath(filepath.Join(env.SSHPath, kt.KeyBaseName)) {
			key.IsDefault = true
		}

		return nil
	})

	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
		return nil
	}

	if key.PublicKey != "" && key.PrivateKey != "" {
		// key desc
		keyDesc := ""
		keyStr := strings.SplitAfterN(getKeyPayload(key.PublicKey), " ", 3)
		if len(keyStr) >= 3 {
			keyDesc = strings.TrimSpace(keyStr[2])
		}
		key.KeyDesc = keyDesc
		return key
	}

	return nil
}

// CreateLink creates symbol link for specified SSH key
func createLink(alias string, keyMap map[string]*SSHKey, env *Environment) error {
	// clears both private & public keys from SSH key path
	key, found := keyMap[alias]
	if !found {
		return errors.New(fmt.Sprintf("ssh key [%s] not found", alias))
	}
	err := clearKey(env)
	if err != nil {
		return err
	}
	// Create symlink for private key
	err = os.Symlink(filepath.Join(env.SKMPath, alias, key.Type.PrivateKey()), filepath.Join(env.SSHPath, key.Type.PrivateKey()))
	if err != nil {
		return errors.Wrap(err, "Create symlink failed")
	}

	// Create symlink for public key
	err = os.Symlink(filepath.Join(env.SKMPath, alias, key.Type.PublicKey()), filepath.Join(env.SSHPath, key.Type.PublicKey()))
	if err != nil {
		return errors.Wrap(err, "Create symlink failed")
	}
	return nil
}

// RunHook runs hook file after switching SSH key
func runHook(alias string, env *Environment) {
	if info, err := os.Stat(filepath.Join(env.StorePath, alias, HookName)); !os.IsNotExist(err) {
		if info.Mode()&0111 != 0 {
			utils.Execute("", filepath.Join(env.SKMPath, alias, HookName), alias)
		}
	}
}

func RenameSSHKey(aliasName string, newaliasName string, env *Environment) error {

	alias := aliasName
	newAlias := newaliasName
	if len(alias) == 0 || len(newAlias) == 0 {
		return errors.New("Please input source alias name and new alias name")
	} else {
		err := os.Rename(filepath.Join(env.SKMPath, alias), filepath.Join(env.SKMPath, newAlias))
		if err == nil {
			utils.PrintN(utils.Info, fmt.Sprintf("SSH key [%s] renamed to [%s]", alias, newAlias))
		} else {
			return errors.New("Failed to rename SSH key!")
		}
	}
	return nil
}

func DeleteSSHKey(aliasName string, env *Environment) error {

	var alias string
	if len(aliasName) > 0 {
		alias = aliasName
	} else {
		return errors.New("Must input SSHKey alias")
	}
	keyMap, _ := LoadSSHKeys(env)
	key, ok := keyMap[alias]
	if !ok {
		return errors.New(fmt.Sprintf("Key alias: %s doesn't exist!", alias))
	}

	if err := deleteKey(alias, key, env); err != nil {
		return errors.Wrap(err, "delete key failed")
	}
	return nil
}

func deleteKey(alias string, key *SSHKey, env *Environment) error {

	var confirm bool
	var err error
	inUse := key.IsDefault
	if inUse {
		c := prompt.NewDefaultConfirm("SSHKey ["+alias+"] is currently in use, please confirm to delete it", false)
		if confirm, err = c.Run(); err != nil {
			return err
		}

	} else {
		c := promptx.NewDefaultConfirm("Please confirm to delete SSHKey ["+alias+"]", true)
		if confirm, err = c.Run(); err != nil {
			return err
		}
	}
	if confirm {
		if inUse {
			err := clearKey(env)
			if err != nil {
				return err
			}
		}
		// Remove specified key by alias name
		err := os.RemoveAll(filepath.Join(env.SKMPath, alias))
		if err != nil {
			return errors.Wrapf(err, "Failed to delete SSH key [%s]!", alias)
		} else {
			utils.PrintN(utils.Warn, fmt.Sprintf("SSH key [%s] deleted!", alias))
		}
	}
	return nil
}

func DisplaySSHKey(aliasName string, env *Environment) error {

	alias := aliasName
	keys, _ := LoadSSHKeys(env)
	if len(alias) > 0 {
		if key, exists := keys[alias]; exists {
			fmt.Print(getKeyPayload(key.PublicKey))
			return nil
		} else {
			return errors.New(fmt.Sprintf("Key alias[%s] not found", alias))
		}
	} else {
		// display default key
		for _, key := range keys {
			if key.IsDefault {
				keyPath := utils.ParseOriginalFilePath(filepath.Join(env.SSHPath, key.Type.PublicKey()))
				fmt.Print(getKeyPayload(keyPath))
				break
			}
		}
	}
	return nil
}

func SetSSHKey(aliasName string, env *Environment) (map[string]*SSHKey, error) {

	alias := aliasName
	keyMap, _ := LoadSSHKeys(env)
	if len(alias) == 0 {

		cfg := &promptx.SelectConfig{
			ActiveTpl:    `» {{ .Alias | cyan }} {{"[" | cyan }}{{.Type.Name | cyan }}{{"]" | cyan }}`,
			InactiveTpl:  `  {{ .Alias | white }} {{"[" | white }}{{.Type.Name | white }}{{"]" | white }}`,
			SelectPrompt: "one ssh key",
			SelectedTpl:  `» {{ .Alias | green | bold }}`,
			DisPlaySize:  9,
			DetailsTpl: `
--------- SSH Key ----------
{{ "Alias:" | blue | faint }} {{ .Alias | blue | faint }}
{{ "Type:" | blue | faint }} {{ .Type.Name | blue | faint }}`,
		}

		// Construct prompt menu items
		var keyList []*SSHKey
		for k := range keyMap {
			keyList = append(keyList, keyMap[k])
		}
		s := &promptx.Select{
			Items:  keyList,
			Config: cfg,
		}
		idx := s.Run()
		alias = keyList[idx].Alias
	}

	// _, ok := keyMap[alias]
	// if !ok {
	// 	return nil, errors.New(fmt.Sprintf("Key alias: %s doesn't exist!", alias))
	// }

	// Set key with related alias as default used key
	err := createLink(alias, keyMap, env)
	if err != nil {
		return nil, err
	}
	// Run a potential hook
	runHook(alias, env)
	utils.PrintN(utils.Info, fmt.Sprintf("Now using SSH key: [%s]", alias))

	sshKey := loadSingleKey(filepath.Join(env.SKMPath, alias), alias, env)
	keys := map[string]*SSHKey{}
	if sshKey != nil {
		keys[alias] = sshKey
	}
	return keys, nil
}

func getKeyPayload(keyPath string) string {
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		fmt.Println("Failed to read ", keyPath)
		return ""
	}
	return string(key)
}

// ClearKey clears both private & public keys from SSH key path
func clearKey(env *Environment) error {

	for _, kt := range SupportedKeyTypes {
		// Remove private key if exists
		PrivateKeyPath := filepath.Join(env.SSHPath, kt.PrivateKey())
		if utils.PathExist(PrivateKeyPath) {
			err := os.Remove(PrivateKeyPath)
			if err != nil {
				return err
			}
		}
		// Remove public key if exists
		PublicKeyPath := filepath.Join(env.SSHPath, kt.PublicKey())
		if utils.PathExist(PublicKeyPath) {
			err := os.Remove(PublicKeyPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// get private key from SSHPath

func getDefaultKey(sshPath string) (error, string) {

	var keyfile string
	keyfile = ""
	for _, kt := range SupportedKeyTypes {
		keyfile = filepath.Join(sshPath, kt.PrivateKey())
		if exist := utils.PathExist(keyfile); exist {
			// check privateKey
			_, err := privateKeyFile(keyfile, "")
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("bad PrivateKeyFile: [%s]", keyfile)), ""
			}
			break
		}
	}

	return nil, keyfile
}
