package myssh

import (
	"fmt"
	"github.com/cnlubo/myssh/utils"
	"github.com/cnlubo/promptx"
	"github.com/fatih/set"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/schwarmco/go-cartesian-product"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type replaceRange struct {
	sourceStr string
	ranges    []string
	rangesit  []interface{}
}

type replaceRanges []replaceRange

var (
	// separator regexp

	baseOperator = regexp.MustCompile(`([-+*/])`)
	operator     = regexp.MustCompile(`[\s]` + baseOperator.String() + `[\s]`)
	allOperator  = regexp.MustCompile(operator.String() + `|(\s+)`)

	comma          = regexp.MustCompile(`[,]`)
	rangeSeparator = regexp.MustCompile(`(?:-|\.\.)`)

	// hostPat regexp

	normalPat = regexp.MustCompile(`^[^\[]+$`)
	rangePat  = regexp.MustCompile(`\w+(?:(?:-|\.\.)\w+)+`)

	// cluster variables
	varNameReg = regexp.MustCompile(`^\w[-.\w]*$`)

	ClustersCfg ClustersConfig
)

func ParseExpr(expr string) ([]string, error) {

	var (
		op          []string
		patList     []string
		hostPatList [][]string
	)
	if ok, errDesc := checkHostRegexp(expr); !ok {
		return []string{""}, errors.New(errDesc)
	}

	hostPattern := utils.DeleteExtraSpace(trimExpr(expr))

	result := allOperator.FindAllStringSubmatch(hostPattern, -1)

	for _, v := range result {
		if len(utils.CompressStr(v[0])) == 0 {
			op = append(op, "+")
		} else {
			for i := 1; i < len(v); i++ {
				if len(v[i]) > 0 {
					op = append(op, utils.CompressStr(v[i]))
				}
			}
		}
	}

	patList = allOperator.Split(hostPattern, -1)
	for _, value := range patList {
		var pats []string
		if reg.MatchString(value) {
			if oneVar.MatchString(value) {
				match := oneVar.FindAllStringSubmatchIndex(value, -1)
				for i := 0; i < len(match); i++ {
					name := value[match[i][2]:match[i][3]]
					if cs, ok := ClustersCfg.Clusters.FindClusterByName(name); ok {
						pats, _ = ParseExpr(cs.HostPattern)

					} else {
						utils.PrintN(utils.Err, fmt.Sprintf("Variable %s not defined.", value))
					}
				}
			} else {
				utils.PrintN(utils.Err, fmt.Sprintf("Invalid variable reference syntax:%s", value))
			}
		} else {
			p := parseHostPattern(value)
			pats = append(pats, p...)
		}

		hostPatList = append(hostPatList, pats)
	}

	setsA := set.New(set.ThreadSafe)
	for _, v := range hostPatList[0] {
		setsA.Add(v)
	}
	for i := 0; i < len(op); i++ {

		setsB := set.New(set.ThreadSafe)
		for _, v := range hostPatList[i+1] {
			setsB.Add(v)
		}
		switch op[i] {

		case string("+"):
			setsA = set.Union(setsA, setsB)

		case string("-"):
			setsA = set.Difference(setsA, setsB)

		case string("*"):
			setsA = set.Intersection(setsA, setsB)

		case string("/"):
			setsA = set.SymmetricDifference(setsA, setsB)
		}

	}

	return set.StringSlice(setsA), nil

}

func AddCluster() (Clusters, error) {

	cluster := ClusterConfig{}
	var err error
	// clusterName
	p := promptx.NewDefaultPrompt(func(line []rune) error {
		if strings.TrimSpace(string(line)) == "" {
			return inputEmptyErr
		} else if len(line) > 12 {
			return inputTooLongErr
		}
		_, ok := ClustersCfg.Clusters.FindClusterByName(string(line))
		if ok {
			return errors.New(fmt.Sprintf("Cluster [%s] already existed", string(line)))
		}
		return nil

	}, "ClusterName:")

	cluster.Name, err = p.Run()
	if err != nil {
		return nil, err
	}
	// HostPatterns
	p = promptx.NewDefaultPrompt(func(line []rune) error {
		if strings.TrimSpace(string(line)) == "" {
			return inputEmptyErr
		}
		if ok, _ := checkHostRegexp(string(line)); !ok {
			return errors.New(fmt.Sprintf("bad host pattern:%s", string(line)))
		}
		return nil

	}, "HostPattern:")

	cluster.HostPattern, err = p.Run()
	if err != nil {
		return nil, err
	}

	// Save
	ClustersCfg.Clusters = append(ClustersCfg.Clusters, &cluster)
	sort.Sort(ClustersCfg.Clusters)

	err = ClustersCfg.write()
	if err != nil {
		return nil, errors.Wrap(err, "add cluster failed")
	}
	fmt.Println()
	utils.PrintN(utils.Info, fmt.Sprintf("Add cluster[%s] success\n", cluster.Name))
	return []*ClusterConfig{&cluster}, nil
}

func DeleteClusters(names []string) (Clusters, error) {

	var deletesIdx []int
	var deleteNodes Clusters
	var deleteNames []string
	for _, deleteName := range names {
		for i, s := range ClustersCfg.Clusters {

			if strings.ToLower(s.Name) == strings.ToLower(deleteName) {
				deletesIdx = append(deletesIdx, i)
				deleteNodes = append(deleteNodes, s)
				deleteNames = append(deleteNames, s.Name)

			}

			// matched, err := filepath.Match(deleteName, s.Name)
			// if err != nil {
			// 	// check equal
			// 	if strings.ToLower(s.Name) == strings.ToLower(deleteName) {
			// 		deletesIdx = append(deletesIdx, i)
			// 	}
			// } else {
			// 	if matched {
			// 		deletesIdx = append(deletesIdx, i)
			// 		deleteNodes = append(deleteNodes, s)
			// 		deleteNames=append(deleteNames,s.Name)
			// 	}
			// }

		}

	}

	var confirm bool
	var err error
	clName := strings.Join(deleteNames, ",")
	if len(deletesIdx) == 0 {
		return nil, errors.New("none clusters delete!!!")
	} else {
		message := "(Please confirm to delete clusters [" + clName + "]"
		c := promptx.NewDefaultConfirm(message,false)
		confirm, err = c.Run()
		if err != nil {
			return nil, err
		}
	}

	if confirm {
		// sort and delete
		sort.Ints(deletesIdx)
		for i, del := range deletesIdx {
			ClustersCfg.Clusters = append(ClustersCfg.Clusters[:del-i], ClustersCfg.Clusters[del-i+1:]...)
		}

		// save config
		sort.Sort(ClustersCfg.Clusters)
		err := ClustersCfg.write()
		if err != nil {
			return nil, errors.Wrap(err, "delete clusters failed")
		}

		utils.PrintN(utils.Info, fmt.Sprintf("delete clusters successfully\n"))
	} else {
		return nil, nil
	}
	return deleteNodes, nil
}

// batch execution of commands
func ClusterBatchCmds(hostPattern string, PromptPass bool, login *ServerConfig, cmd ...string) error {

	hostList, err := ParseExpr(hostPattern)
	if err != nil {
		return errors.Wrap(err, "bad hostPattern")
	}
	var servers Servers

	reg := regexp.MustCompile(`^(\w+)[@](.+$)`) // user@host:port
	for _, value := range hostList {
		var sshUser, loginPass, host, port, privateKey, authMethod string

		if reg.MatchString(utils.CompressStr(value)) {
			sshUser, host, port = utils.ParseConnect(utils.CompressStr(value))

		} else {
			sshUser = login.User
			host = utils.CompressStr(value)
			port = strconv.Itoa(login.Port)
		}
		pp, _ := strconv.Atoi(port)

		// set default
		if sshUser == "" {
			sshUser = ClustersCfg.Default.User
		}
		if pp == 0 {
			pp = ClustersCfg.Default.Port
		}
		if login.PrivateKey == "" {
			privateKey = ClustersCfg.Default.PrivateKey
		} else {
			privateKey = login.PrivateKey
		}
		// prompt pass
		if PromptPass {
			utils.PrintN(utils.Warn, "Server "+sshUser+"@"+host)
			loginPass = ""

			prompt := promptx.NewDefaultPrompt(func(line []rune) error {
				if strings.TrimSpace(string(line)) == "" {
					return inputEmptyErr
				}
				return nil

			}, "Please type login password:")

			prompt.Mask = MaskPrompt
			loginPass, err = prompt.Run()
			if err != nil {
				return err
			}
			authMethod = "password"

		} else {
			loginPass = ""
			authMethod = "key"
		}
		server := &ServerConfig{
			Name:       host,
			User:       sshUser,
			Address:    host,
			Port:       pp,
			Method:     authMethod,
			PrivateKey: privateKey,
			// PrivateKeyPassword: login.PrivateKeyPassword,
			Password:            loginPass,
			ServerAliveInterval: ClustersCfg.Default.ServerAliveInterval,
		}
		servers = append(servers, server)

	}
	exCmd := strings.Join(cmd, "&&")
	err = batchExec(servers, exCmd)
	if err != nil {
		return err
	}
	return nil
}

func ClusterKeyCopy(hostPattern string, sshPort int, sshUser string, identityfile string) error {

	hostList, err := ParseExpr(hostPattern)
	if err != nil {
		return errors.Wrap(err, "bad hostPattern")
	}
	// set default
	if sshPort == 0 {
		sshPort = ClustersCfg.Default.Port
	}
	if len(sshUser) == 0 {
		sshUser = ClustersCfg.Default.User
	}

	if len(utils.CompressStr(identityfile)) == 0 {
		identityfile = ClustersCfg.Default.PrivateKey
	}

	// check identityfile

	if len(identityfile) > 0 {
		home, _ := homedir.Dir()
		identityfile = utils.ParseRelPath(identityfile, home+"/.ssh")

		_, err := privateKeyFile(identityfile, " ")
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("cluster (%s) have bad identityfile", hostPattern))
		}
	} else {
		return errors.New(fmt.Sprintf("cluster (%s) have empty identityfile", hostPattern))
	}

	for _, host := range hostList {
		connectStr := sshUser + "@" + host
		utils.PrintN(utils.Info, fmt.Sprintf("copy current SSH key to [%s]\n", connectStr))
		result := CopyKey(connectStr, strconv.Itoa(sshPort), identityfile)
		if !result {
			utils.PrintN(utils.Err, fmt.Sprintf("Copy SSHKey to [%s]failed", connectStr))
		}
	}

	return nil
}

func checkHostRegexp(exp string) (bool, string) {

	var result = true
	var errDesc string
	ex := utils.CompressStr(exp)
	// Expecting terms but found operator
	reg := regexp.MustCompile(`^\s*` + baseOperator.String() + `|\s+` + baseOperator.String() + `\s*$`)
	if reg.MatchString(ex) {
		s := reg.FindAllStringSubmatch(ex, -1)
		var errOp string

		for i := 0; i < len(s); i++ {
			if len(errOp) == 0 {
				errOp = errOp + utils.CompressStr(s[i][0])
			} else {
				errOp = errOp + " " + utils.CompressStr(s[i][0])
			}
		}
		errDesc = fmt.Sprintf("Expecting terms but found operator:[%s]", errOp) + "\n"

		result = false
	}

	// Expecting operators but found term
	reg = regexp.MustCompile(`(?:}\w*)({.*?})`)

	if reg.MatchString(ex) {
		s := reg.FindAllStringSubmatch(ex, -1)
		var errTerm string
		for i := 0; i < len(s); i++ {
			if len(errTerm) == 0 {
				errTerm = errTerm + utils.CompressStr(s[i][1])
			} else {
				errTerm = errTerm + " " + utils.CompressStr(s[i][1])
			}
		}
		if len(errDesc) == 0 {
			errDesc = fmt.Sprintf("Expecting operators but found term:[%s]", errTerm) + "\n"
		} else {
			errDesc = errDesc + fmt.Sprintf("Expecting operators but found term:[%s]", errTerm) + "\n"
		}
		result = false
	}
	// check variable
	reg = regexp.MustCompile(`{(.*?)}`)
	if reg.MatchString(ex) {
		match := reg.FindAllStringSubmatchIndex(ex, -1)
		var errName string
		for i := 0; i < len(match); i++ {
			name := ex[match[i][2]:match[i][3]]
			if ok := varNameReg.MatchString(name); !ok {
				if len(errName) == 0 {
					errName = errName + ex[match[i][0]:match[i][1]]
				} else {
					errName = errName + " " + ex[match[i][0]:match[i][1]]
				}
			}
		}
		if len(errName) > 0 {
			if len(errDesc) == 0 {
				errDesc = fmt.Sprintf("Invalid variable name: %s", errName) + "\n"
			} else {
				errDesc = errDesc + fmt.Sprintf("Invalid variable name: %s", errName) + "\n"
			}
			result = false
		}
	}

	// check variable reference syntax
	reg = regexp.MustCompile(`[})]([-+*/])[^{(\s]|[^})\s]([-+*/])[{(]`)
	if reg.MatchString(ex) {
		s := reg.FindAllStringSubmatch(ex, -1)
		fmt.Println(s)
		var errSyntax string
		for i := 0; i < len(s); i++ {
			if len(errSyntax) == 0 {
				errSyntax = errSyntax + utils.CompressStr(s[i][0])
			} else {
				errSyntax = errSyntax + " " + utils.CompressStr(s[i][0])
			}
		}
		if len(errDesc) == 0 {
			errDesc = fmt.Sprintf("Invalid variable reference syntax:[%s]", errSyntax) + "\n"
		} else {
			errDesc = errDesc + fmt.Sprintf("Invalid variable reference syntax:[%s]", errSyntax) + "\n"
		}
		result = false
	}

	return result, errDesc
}

func replaceAllSubmatchFunc(re *regexp.Regexp, content []byte, f func(s []byte) string) []byte {

	indexes := re.FindAllSubmatchIndex(content, -1)

	if len(indexes) == 0 {
		return content
	}

	var idxes [][]int
	for _, index := range indexes {
		var v []int
		for j := 0; j < len(index); j++ {
			if index[j] != -1 {
				v = append(v, index[j])
			}
		}
		idxes = append(idxes, v)
	}

	l := len(idxes)
	ret := append([]byte{}, content[:idxes[0][2]]...)

	for i, pair := range idxes {
		switch {
		case pair[0] == pair[2]:
			s := strings.TrimLeft(f(content[pair[2]:pair[3]]), " ")
			ret = append(ret, []byte(s)...)
		case pair[1] == pair[3]:
			s := strings.TrimRight(f(content[pair[2]:pair[3]]), " ")
			ret = append(ret, []byte(s)...)
		default:
			s := f(content[pair[2]:pair[3]])
			ret = append(ret, []byte(s)...)
		}

		if i+1 < l {

			switch {
			case pair[1] == pair[3]:
				if idxes[i+1][0] == idxes[i+1][2] {
					ret = append(ret, content[pair[1]:idxes[i+1][0]]...)
				} else {
					ret = append(ret, content[pair[1]:idxes[i+1][0]+1]...)
				}
			default:
				if idxes[i+1][0] == idxes[i+1][2] {
					ret = append(ret, content[pair[1]-1:idxes[i+1][0]]...)
				} else {
					ret = append(ret, content[pair[1]-1:idxes[i+1][0]+1]...)
				}
			}
		}
	}

	if idxes[len(idxes)-1][1] == idxes[len(idxes)-1][3] {
		ret = append(ret, content[idxes[len(idxes)-1][1]:]...)
	} else {
		ret = append(ret, content[idxes[len(idxes)-1][1]-1:]...)
	}
	return ret
}

func trimExpr(expr string) string {

	var opRegexp = regexp.MustCompile(`[})]([-+*/])[{(]?|[})]?([-+*/])[{(]`)
	ex := utils.CompressStr(expr)
	if ok := opRegexp.MatchString(ex); ok {
		result := replaceAllSubmatchFunc(opRegexp, []byte(ex), func(s []byte) string {
			m := " " + string(s) + " "
			return m
		})
		ex = string(result)
	}
	return ex

}

var (
	reg    = regexp.MustCompile(`{.*?}`)
	oneVar = regexp.MustCompile(`^{([^{}\\\s]*)}$`)
)

var variable = regexp.MustCompile(`\[([^]]*)]`)

func parseHostPattern(pattern string) []string {

	var hostList []string
	if len(pattern) == 0 {
		return nil
	}

	pat := utils.CompressStr(pattern)

	if variable.MatchString(pat) {

		variableIndexes := variable.FindAllStringIndex(pat, -1)

		if r, _ := parseVariables(pat, variableIndexes); r != nil {

			var reps chan []interface{}

			switch {

			case len(variableIndexes) == 1:
				reps = cartesian.Iter(r[0].rangesit)
			case len(variableIndexes) == 2:

				reps = cartesian.Iter(r[0].rangesit, r[1].rangesit)

			case len(variableIndexes) == 3:

				reps = cartesian.Iter(r[0].rangesit, r[1].rangesit, r[2].rangesit)

			case len(variableIndexes) == 4:

				reps = cartesian.Iter(r[0].rangesit, r[1].rangesit, r[2].rangesit, r[3].rangesit)

			}

			if len(variableIndexes) <= 4 {

				for cs := range reps {
					result := pat
					for k := range r {
						result = strings.Replace(result, r[k].sourceStr, cs[k].(string), 1)

					}

					hostList = append(hostList, result)
				}
			} else {
				utils.PrintN(utils.Err, fmt.Sprintf("too many range:%s", pattern))
			}
		}
	} else {
		if normalPat.MatchString(pat) {
			s := normalPat.FindStringSubmatch(pat)[0]
			hostList = append(hostList, s)

		} else {
			utils.PrintN(utils.Err, fmt.Sprintf("bad host:%s", pattern))

		}
	}

	return hostList
}

func endNum(s string) string {
	var result string
	if utils.IsNumeric(s) {
		result = s
	} else {

		for i := 0; i < len(s); i++ {
			result += "9"

		}

	}
	return result
}

func parseRanges(variables []string) ([]string, error) {

	if len(variables) == 0 {
		return nil, nil
	}

	var rangeVariable []string

	for _, v := range variables {

		rr := rangeSeparator.Split(v, -1)
		var bs, es int

		if utils.IsNumeric(rr[0]) {
			bs, _ = strconv.Atoi(rr[0])
			es, _ = strconv.Atoi(endNum(rr[1]))

			if bs > es {
				utils.PrintN(utils.Err, fmt.Sprintf("bad range:%s", v))
			} else {

				for i := bs; i <= es; i++ {

					rangeVariable = append(rangeVariable, strconv.Itoa(i))

				}
			}
		} else {
			bs := strings.Index(utils.Alphabet, rr[0])
			es := strings.LastIndex(utils.Alphabet, rr[1])
			for i := bs; i <= es; i++ {

				rangeVariable = append(rangeVariable, string(utils.Alphabet[i]))

			}

		}

	}

	return rangeVariable, nil
}

var (
	// valid regexp
	validRangePat = regexp.MustCompile(`\w+(?:(?:-|\.\.)\w+)?`)
	validVariable = regexp.MustCompile(`^` + validRangePat.String() + `(?:\s*,\s*` + validRangePat.String() + `)*$`)
)

func parseVariables(hostPattern string, indexes [][]int) (replaceRanges, error) {

	var re replaceRanges

	for _, r := range indexes {

		s := hostPattern[r[0]+1 : r[1]-1]

		if ok := validVariable.MatchString(s); ok {

			variables := comma.Split(hostPattern[r[0]+1:r[1]-1], -1)

			var ranges []string
			var vars []string

			for _, v := range variables {

				if rangePat.MatchString(v) {
					ranges = append(ranges, v)
				} else {
					vars = append(vars, v)
				}
			}

			if rr, _ := parseRanges(ranges); rr != nil {
				vars = utils.MergeSlice(vars, rr)
			}

			it := make([]interface{}, len(vars))

			for i, v := range vars {
				it[i] = v
			}
			r := replaceRange{
				sourceStr: hostPattern[r[0]:r[1]],
				ranges:    vars,
				rangesit:  it,
			}
			re = append(re, r)

		} else {

			utils.PrintN(utils.Err, fmt.Sprintf("Bad range:%s", s))

		}

	}

	if re != nil {
		return re, nil
	} else {
		return nil, nil
	}
}

func ClustersConfigExample(env *Environment) *ClustersConfig {
	return &ClustersConfig{
		Default:  defaultClusterExample(env),
		Clusters: Clusters{},
	}
}
func defaultClusterExample(env *Environment) DefaultClusterConfig {
	_, privateKey := getDefaultKey(env.SSHPath)

	return DefaultClusterConfig{
		User:                utils.GetUsername(),
		Port:                22,
		PrivateKey:          privateKey,
		ServerAliveInterval: 30 * time.Second,
	}
}

// write clustersConfigfile
func (cfg *ClustersConfig) write() error {
	if cfg.configPath == "" {
		return errors.New("nodes config path not set")
	}
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cfg.configPath, out, 0644)
}

// write nodes config to yaml file
func (cfg *ClustersConfig) WriteTo(filePath string) error {
	if filePath == "" {
		return errors.New("file path is empty")
	}
	cfg.configPath = filePath
	return cfg.write()
}

// load nodes config
func (cfg *ClustersConfig) load() error {
	if cfg.configPath == "" {
		return errors.New("config path not set")
	}
	buf, err := ioutil.ReadFile(cfg.configPath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(buf, cfg)
}

// load nodes from yaml file
func (cfg *ClustersConfig) LoadFrom(filePath string) error {
	if filePath == "" {
		return errors.New("file path is empty")
	}
	cfg.configPath = filePath
	return cfg.load()
}

// find cluster by name
func (cs Clusters) FindClusterByName(name string) (*ClusterConfig, bool) {
	for _, c := range cs {
		if name == c.Name {
			return c, true
		}
	}
	return nil, false
}
