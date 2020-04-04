package main

import (
	"bufio"
	"bytes"
	b64 "encoding/base64"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

/*

curl -H "Private-Token: NAq7xdBggmbkNLkpsQm6" https://pnl.gitlab.schubergphilis.com/api/v4/projects/



ROUTINE NOW:

1. Paar keer (5) via Curl de projecten ophalen uit Gitlab

curl "https://pnl.gitlab.schubergphilis.com/api/v4/projects?private_token=NAq7xdBggmbkNLkpsQm6?page=1&per_page=100" > output1.json

Hier wordt pagination gebruikt.

2. Met go functie ReadProjectFile wordt een paar keer een CSV gemaakt aan de hand van het aantal files uit [1].

3. De CSVs uit [2] worden gecombineerd tot één grote CSV.

4. Per projectID uit lijst-3, haal de base64 encoded pom.xml filecontent op, decode dit terug naar string en schrijf dit naar een nieuw bestand samen met projectnaam.

6. En daarna iets met XML parsing van de SharedModules uit het bestand en dit ergens naartoe printen zodat per project te zien is welke sharedmodules er in de POM file stonden.

- - - - - - - - - - -


- - - - - - - - - - -

ROUTINE TO BE:

1. Crawl projectlijst van SBP Statuspage
2. Per item van lijst-1, zoek het project op in Gitlab en haal het projectID op en schrijf dit naar een bestand samen met het de projectnaam.
3. Per projectID uit lijst-2, haal de base64 encoded pom.xml filecontent op, decode dit terug naar string en schrijf dit naar een nieuw bestand samen met projectnaam.
4. En daarna iets met XML parsing van de SharedModules uit het bestand en dit ergens naartoe printen zodat per project te zien is welke sharedmodules er in de POM file stonden.

- - - - - - - - - - -

LIJST COMPONENTEN OPHALEN VAN DE STATUSPAGE
Crawl de webnpagina

LIJST MET GITLAB PROJECT OPHALEN AAN DE HAND VAN DE LIJST COMPONENTEN VAN DE STATUSPAGE
curl -H "Private-Token: NAq7xdBggmbkNLkpsQm6" https://pnl.gitlab.schubergphilis.com/api/v4/projects/824

OPHALEN BASE64 ENCODED BESTAND EN WEGSCHRIJVEN NAAR BESTANDEN
curl --request GET --header 'PRIVATE-TOKEN: NAq7xdBggmbkNLkpsQm6' 'https://pnl.gitlab.schubergphilis.com/api/v4/projects/796/repository/files/tibco%2Fpom%2Exml?ref=master'

curl 'https://pnl.gitlab.schubergphilis.com/api/v4/projects/796/repository/files/tibco%2Fpom%2Exml?ref=master&private_token=NAq7xdBggmbkNLkpsQm6'

ref: https://docs.gitlab.com/ee/api/repository_files.html
ref: https://stackoverflow.com/questions/44730632/gitlab-api-how-to-get-the-repository-project-files-and-metadata


*/

//go run.gitlab.go > /Users/joshuamoesa/Desktop/Go-ModuleInfo/modules_sharedmodules.txt

//Jayway JsonPath Evaluator: https://jsonpath.herokuapp.com/

//Reference structs generator: https://mholt.github.io/json-to-go/
//Reference structs: https://gobyexample.com/structs

//Gitlab file struct
type AutoGenerated struct {
	FileName      string `json:"file_name"`
	FilePath      string `json:"file_path"`
	Size          int    `json:"size"`
	Encoding      string `json:"encoding"`
	ContentSha256 string `json:"content_sha256"`
	Ref           string `json:"ref"`
	BlobID        string `json:"blob_id"`
	CommitID      string `json:"commit_id"`
	LastCommitID  string `json:"last_commit_id"`
	Content       string `json:"content"`
}

//Project struct
type Project []struct {
	ID                int           `json:"id"`
	Description       string        `json:"description"`
	Name              string        `json:"name"`
	NameWithNamespace string        `json:"name_with_namespace"`
	Path              string        `json:"path"`
	PathWithNamespace string        `json:"path_with_namespace"`
	CreatedAt         time.Time     `json:"created_at"`
	DefaultBranch     string        `json:"default_branch"`
	TagList           []interface{} `json:"tag_list"`
	SSHURLToRepo      string        `json:"ssh_url_to_repo"`
	HTTPURLToRepo     string        `json:"http_url_to_repo"`
	WebURL            string        `json:"web_url"`
	ReadmeURL         string        `json:"readme_url"`
	AvatarURL         interface{}   `json:"avatar_url"`
	StarCount         int           `json:"star_count"`
	ForksCount        int           `json:"forks_count"`
	LastActivityAt    time.Time     `json:"last_activity_at"`
	Namespace         struct {
		ID        int         `json:"id"`
		Name      string      `json:"name"`
		Path      string      `json:"path"`
		Kind      string      `json:"kind"`
		FullPath  string      `json:"full_path"`
		ParentID  int         `json:"parent_id"`
		AvatarURL interface{} `json:"avatar_url"`
		WebURL    string      `json:"web_url"`
	} `json:"namespace"`
	Links struct {
		Self          string `json:"self"`
		Issues        string `json:"issues"`
		MergeRequests string `json:"merge_requests"`
		RepoBranches  string `json:"repo_branches"`
		Labels        string `json:"labels"`
		Events        string `json:"events"`
		Members       string `json:"members"`
	} `json:"_links"`
	Archived                                  bool          `json:"archived"`
	Visibility                                string        `json:"visibility"`
	ResolveOutdatedDiffDiscussions            bool          `json:"resolve_outdated_diff_discussions"`
	ContainerRegistryEnabled                  bool          `json:"container_registry_enabled"`
	IssuesEnabled                             bool          `json:"issues_enabled"`
	MergeRequestsEnabled                      bool          `json:"merge_requests_enabled"`
	WikiEnabled                               bool          `json:"wiki_enabled"`
	JobsEnabled                               bool          `json:"jobs_enabled"`
	SnippetsEnabled                           bool          `json:"snippets_enabled"`
	SharedRunnersEnabled                      bool          `json:"shared_runners_enabled"`
	LfsEnabled                                bool          `json:"lfs_enabled"`
	CreatorID                                 int           `json:"creator_id"`
	ImportStatus                              string        `json:"import_status"`
	OpenIssuesCount                           int           `json:"open_issues_count"`
	PublicJobs                                bool          `json:"public_jobs"`
	CiConfigPath                              interface{}   `json:"ci_config_path"`
	SharedWithGroups                          []interface{} `json:"shared_with_groups"`
	OnlyAllowMergeIfPipelineSucceeds          bool          `json:"only_allow_merge_if_pipeline_succeeds"`
	RequestAccessEnabled                      bool          `json:"request_access_enabled"`
	OnlyAllowMergeIfAllDiscussionsAreResolved bool          `json:"only_allow_merge_if_all_discussions_are_resolved"`
	PrintingMergeRequestLinkEnabled           bool          `json:"printing_merge_request_link_enabled"`
	MergeMethod                               string        `json:"merge_method"`
	ExternalAuthorizationClassificationLabel  interface{}   `json:"external_authorization_classification_label"`
	Permissions                               struct {
		ProjectAccess interface{} `json:"project_access"`
		GroupAccess   struct {
			AccessLevel       int `json:"access_level"`
			NotificationLevel int `json:"notification_level"`
		} `json:"group_access"`
	} `json:"permissions"`
	ApprovalsBeforeMerge int  `json:"approvals_before_merge"`
	Mirror               bool `json:"mirror"`
	PackagesEnabled      bool `json:"packages_enabled"`
}

//GitlabProject struct
type GitlabProject struct {
	ID   string
	Name string
}

//PomProject struct
type PomProject struct {
	XMLName        xml.Name `xml:"project"`
	Text           string   `xml:",chardata"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Xmlns          string   `xml:"xmlns,attr"`
	Xsi            string   `xml:"xsi,attr"`
	ModelVersion   string   `xml:"modelVersion"`
	GroupId        string   `xml:"groupId"`
	ArtifactId     string   `xml:"artifactId"`
	Version        string   `xml:"version"`
	Packaging      string   `xml:"packaging"`
	Properties     struct {
		Text                     string `xml:",chardata"`
		MavenPluginVersion       string `xml:"maven_plugin_version"`
		BwceLoggingVersion       string `xml:"bwce.logging.version"`
		BwceMessagebrokerVersion string `xml:"bwce.messagebroker.version"`
	} `xml:"properties"`
	Modules struct {
		Text   string `xml:",chardata"`
		Module string `xml:"module"`
	} `xml:"modules"`
	DistributionManagement struct {
		Text       string `xml:",chardata"`
		Repository struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id"`
			Name string `xml:"name"`
			URL  string `xml:"url"`
		} `xml:"repository"`
		SnapshotRepository struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id"`
			Name string `xml:"name"`
			URL  string `xml:"url"`
		} `xml:"snapshotRepository"`
	} `xml:"distributionManagement"`
	Repositories struct {
		Text       string `xml:",chardata"`
		Repository []struct {
			Text     string `xml:",chardata"`
			ID       string `xml:"id"`
			Name     string `xml:"name"`
			URL      string `xml:"url"`
			Releases struct {
				Text    string `xml:",chardata"`
				Enabled string `xml:"enabled"`
			} `xml:"releases"`
			Snapshots struct {
				Text    string `xml:",chardata"`
				Enabled string `xml:"enabled"`
			} `xml:"snapshots"`
		} `xml:"repository"`
	} `xml:"repositories"`
	Build struct {
		Text       string `xml:",chardata"`
		Extensions struct {
			Text      string `xml:",chardata"`
			Extension struct {
				Text       string `xml:",chardata"`
				GroupId    string `xml:"groupId"`
				ArtifactId string `xml:"artifactId"`
				Version    string `xml:"version"`
			} `xml:"extension"`
		} `xml:"extensions"`
	} `xml:"build"`
}

//XMLQuery struct
type XMLQuery struct {
	Loc string `xml:",chardata"`
}

var xmlquery XMLQuery

type Node struct {
	XMLName xml.Name
	Content []byte `xml:",innerxml"`
	Nodes   []Node `xml:",any"`
}

func perror(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	//readProjectFile("/Users/joshuamoesa/Desktop/Go-ModuleInfo/listprojects5.json")

	//readGitlabProjectFile("/Users/joshuamoesa/Desktop/Go-ModuleInfo/GitlabProjectsByIdTest.csv")
	//readGitlabProjectFile("/Users/joshuamoesa/Desktop/Go-ModuleInfo/GitlabProjectsById.csv")
	//createReport("/Users/joshuamoesa/Desktop/Go-ModuleInfo/pom/")
	//createReport2("/Users/joshuamoesa/Desktop/Go-ModuleInfo/pom/pom_esb-ces-emp-processexecutionplan.xml")
	//createReport3("/Users/joshuamoesa/Desktop/Go-ModuleInfo/pom/")
	//createReport3Archive("/Users/joshuamoesa/Desktop/Go-ModuleInfo/pom/pom_esb-ces-emp-processexecutionplan.xml")

	exampleReadDir()
}

func readRemotePomFile() {

	////https://stackoverflow.com/questions/17156371/how-to-get-json-response-in-golang

	url := "https://pnl.gitlab.schubergphilis.com/api/v4/projects/796/repository/files/tibco%2Fpom%2Exml?ref=master&private_token=NAq7xdBggmbkNLkpsQm6"

	res, err := http.Get(url)

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err.Error())
	}

	AutoGenerated1 := AutoGenerated{}

	json.Unmarshal(body, &AutoGenerated1)
	fmt.Println(AutoGenerated1.Content)

	// var data AutoGenerated
	// json.Unmarshal(body, &data)
	// fmt.Println()

	// os.Exit(0)
}

func readProjectFile(fileName string) {

	//Reference: Schrijven naar bestanden: https://golangbot.com/write-files/

	//fmt.Println("read project file")

	jsonFile, err := os.Open(fileName)

	if err != nil {
		panic(err.Error())
	}

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	//fmt.Println("file read")

	// we initialize our Projects array
	Projects := Project{}

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'projects' which we defined above
	json.Unmarshal(byteValue, &Projects)

	//Create a file
	now := time.Now()
	sec := now.Unix()
	f, err := os.Create("/Users/joshuamoesa/Desktop/Go-ModuleInfo/readProjectFileOutput" + strconv.FormatInt(sec, 10) + ".csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	for i := 0; i < len(Projects); i++ {
		ProjectID := strconv.Itoa(Projects[i].ID)
		// fmt.Println("Project id: " + ProjectID)
		// fmt.Println("Project name: " + Projects[i].Name)
		fmt.Println(ProjectID + Projects[i].Name)

		l, err := f.WriteString(ProjectID + "," + Projects[i].Name + "\n")
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
		fmt.Println(l, "bytes written successfully")

	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

}

func readGitlabProjectFile(fileName string) {

	//Reference: omgaan met CSV 	https://www.thepolyglotdeveloper.com/2017/03/parse-csv-data-go-programming-language/

	csvFile, _ := os.Open(fileName)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	//var project []GitlabProject

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		fmt.Println("Working for " + line[1])

		//Reference: https://docs.gitlab.com/ce/api/repository_files.html#get-raw-file-from-repository
		url := "https://pnl.gitlab.schubergphilis.com/api/v4/projects/" + line[0] + "/repository/files/tibco%2Fpom%2Exml?ref=tst&private_token=NAq7xdBggmbkNLkpsQm6"
		fmt.Println("Search url:" + url)

		res, err := http.Get(url)

		if err != nil {
			panic(err.Error())
		}

		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			panic(err.Error())
		}

		AutoGenerated1 := AutoGenerated{}

		// fetch JSON from an API: https://blog.alexellis.io/golang-json-api-client/
		json.Unmarshal(body, &AutoGenerated1)

		//Reference: Golang decode base64 data to string https://stackoverflow.com/questions/46669782/golang-decode-base64-data-to-string

		sDec, _ := b64.StdEncoding.DecodeString(AutoGenerated1.Content)

		//fmt.Println(AutoGenerated1.Content)

		//Create a file
		f, err := os.Create("/Users/joshuamoesa/Desktop/Go-ModuleInfo/pom_" + line[1] + ".xml")
		if err != nil {
			fmt.Println(err)
			return
		}

		l, err := f.WriteString(string(sDec))
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
		fmt.Println(l, "bytes written successfully")

		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

	}

}

func createReport(filePath string) {
	var files []string

	root := filePath
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {

		fmt.Println("Processing: " + file)

		// data, _ := ioutil.ReadFile(file)

		// note := &PomProject{}

		// _ = xml.Unmarshal([]byte(data), &note)

		// fmt.Println(note.Modules.Module)
		// fmt.Println(note.Properties)

	}
}

func createReport2(filePath string) {

	//Reference: https://www.socketloop.com/tutorials/golang-read-xml-elements-data-with-xml-chardata-example
	//

	b, err := ioutil.ReadFile(filePath) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	XMLdata := string(b)

	//Reference: https://stackoverflow.com/questions/17156371/how-to-get-json-response-in-golang

	// example on handling XML chardata(string)
	decoder := xml.NewDecoder(strings.NewReader(string(XMLdata)))

	for {

		// err is ignore here. IF you are reading from a XML file
		// do not ignore err and also check for io.EOF
		token, _ := decoder.Token()

		if token == nil {
			break
		}

		switch Element := token.(type) {
		case xml.StartElement:
			if Element.Name.Local == "properties" {
				fmt.Println("Element name is : ", Element.Name.Local)

				err := decoder.DecodeElement(&xmlquery, &Element)
				if err != nil {
					fmt.Println(err)
				}

				fmt.Println("Element value is : ", xmlquery.Loc)
			}

		// print out the element data
		// convert to []byte slice and cast to string type

		case xml.CharData:
			// str := string([]byte(Element))
			// fmt.Println(str)
		}
	}
}

func createReport3(filePath string) {

	//Reference: https://stackoverflow.com/questions/30256729/how-to-traverse-through-xml-data-in-golang
	//Sample: https://play.golang.org/p/rv1LlxaHvK

	//References:
	// https://stackoverflow.com/questions/12398925/go-xml-marshalling-and-the-root-element
	// https://stackoverflow.com/questions/30256729/how-to-traverse-through-xml-data-in-golang
	// https://www.socketloop.com/tutorials/golang-read-xml-elements-data-with-xml-chardata-example
	// https://tutorialedge.net/golang/parsing-xml-with-golang/

	var files []string

	root := filePath
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, file := range files {

		//fmt.Println("Processing: " + file)

		b, err1 := ioutil.ReadFile(file) // just pass the file name
		if err1 != nil {
			fmt.Print(err1)
		}

		buf := bytes.NewBuffer(b)
		dec := xml.NewDecoder(buf)

		var n Node

		err2 := dec.Decode(&n)
		if err2 != nil {
			panic(err2)
		}

		// walk([]Node{n}, func(n Node) bool {

		// 	if n.XMLName.Local == "artifactId" {
		// 		fmt.Println(string(n.Content))
		// 	}
		// 	if n.XMLName.Local == "properties" {
		// 		fmt.Println(string(n.Content))
		// 	}

		// 	return true
		// })

	}
}

func walk(nodes []Node, f func(Node) bool) {
	for _, n := range nodes {
		if f(n) {
			walk(n.Nodes, f)
		}
	}
}

func exampleReadDir() {

	//http://www.golangprograms.com/example-readall-readdir-and-readfile-from-io-package.html
	//List the files in a folder: https://flaviocopes.com/go-list-files/
	//Reference: https://stackoverflow.com/questions/30256729/how-to-traverse-through-xml-data-in-golang

	entries, err := ioutil.ReadDir("/Users/joshuamoesa/Desktop/Go-ModuleInfo/pom/")
	if err != nil {
		log.Panicf("failed reading directory: %s", err)
	}
	//fmt.Printf("\nNumber of files in current directory: %d", len(entries))
	//fmt.Printf("\nError: %v", err)

	for _, file := range entries {
		//fmt.Println(file.Name())

		//fmt.Println("Processing: " + file)

		b, err1 := ioutil.ReadFile("/Users/joshuamoesa/Desktop/Go-ModuleInfo/pom/" + file.Name()) // just pass the file name
		if err1 != nil {
			fmt.Print(err1)
		}

		buf := bytes.NewBuffer(b)
		dec := xml.NewDecoder(buf)

		var n Node

		err2 := dec.Decode(&n)
		if err2 != nil {
			panic(err2)
		}

		walk([]Node{n}, func(n Node) bool {

			if n.XMLName.Local == "artifactId" {
				fmt.Println(string(n.Content))
			}
			if n.XMLName.Local == "properties" && len(n.Content) > 0 {
				fmt.Println(string(n.Content))
			}

			return true
		})

	}

}

func createReport3Archive(filePath string) {

	//Reference: https://stackoverflow.com/questions/30256729/how-to-traverse-through-xml-data-in-golang
	//Sample: https://play.golang.org/p/rv1LlxaHvK

	b, err := ioutil.ReadFile(filePath) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	buf := bytes.NewBuffer(b)
	dec := xml.NewDecoder(buf)

	var n Node
	err2 := dec.Decode(&n)
	if err2 != nil {
		panic(err2)
	}

	walk([]Node{n}, func(n Node) bool {

		if n.XMLName.Local == "artifactId" {
			artifactIdVar := string(n.Content)
			//fmt.Println(string(n.Content))
			fmt.Println(artifactIdVar)
		}
		if n.XMLName.Local == "properties" {
			propertiesVar := string(n.Content)
			//fmt.Println(string(n.Content))
			fmt.Println(propertiesVar)
		}

		return true
	})

}

func archive() {

	//https://www.socketloop.com/tutorials/golang-convert-http-response-body-to-string

	resp, err := http.Get("https://pnl.gitlab.schubergphilis.com/api/v4/projects/796/repository/files/tibco%2Fpom%2Exml?ref=master&private_token=NAq7xdBggmbkNLkpsQm6")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	htmlData, err := ioutil.ReadAll(resp.Body) //<--- here!

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// print out
	fmt.Println(os.Stdout, string(htmlData)) //<-- here !

	// use Regular Expression to search for keyword
	// for example
	//verified, err := regexp.MatchString("VERIFIED", string(htmlData))

	//if err != nil {
	//		fmt.Println(err)
	//		return
	//	}

}
