// Copyright 2015 ThoughtWorks, Inc.

// This file is part of getgauge/html-report.

// getgauge/html-report is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// getgauge/html-report is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with getgauge/html-report.  If not, see <http://www.gnu.org/licenses/>.
package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/getgauge/common"
	gm "github.com/getgauge/html-report/gauge_messages"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

type summary struct {
	Total   int
	Failed  int
	Passed  int
	Skipped int
}

type overview struct {
	ProjectName string
	Env         string
	Tags        string
	SuccRate    float32
	ExecTime    string
	Timestamp   string
	Summary     *summary
	BasePath    string
}

type specsMeta struct {
	SpecName   string
	ExecTime   string
	Failed     bool
	Skipped    bool
	Tags       []string
	ReportFile string
}

type sidebar struct {
	IsBeforeHookFailure bool
	Specs               []*specsMeta
}

type hookFailure struct {
	HookName   string
	ErrMsg     string
	Screenshot string
	StackTrace string
}

type specHeader struct {
	SpecName string
	ExecTime string
	FileName string
	Tags     []string
	Summary  *summary
}

type row struct {
	Cells []string
	Res   status
}

type table struct {
	Headers []string
	Rows    []*row
}

type spec struct {
	CommentsBeforeTable []string
	Table               *table
	CommentsAfterTable  []string
	Scenarios           []*scenario
	BeforeHookFailure   *hookFailure
	AfterHookFailure    *hookFailure
	Errors              []error
}

type errorType int

type buildError struct {
	ErrorType  errorType
	FileName   string
	LineNumber int
	Message    string
}

func (e buildError) Error() string {
	if e.isParseError() {
		return "[Parse Error] " + e.Message
	}
	return "[Validation Error] " + e.Message
}

func (e buildError) isParseError() bool {
	return e.ErrorType == parseError
}

const (
	parseError errorType = iota
	validationError
)

type scenario struct {
	Heading           string
	ExecTime          string
	Tags              []string
	ExecStatus        status
	Contexts          []item
	Items             []item
	Teardown          []item
	BeforeHookFailure *hookFailure
	AfterHookFailure  *hookFailure
	TableRowIndex     int
}

const (
	stepKind kind = iota
	commentKind
	conceptKind
)

type kind int

type item interface {
	kind() kind
}

type step struct {
	Fragments       []*fragment
	Res             *result
	PreHookFailure  *hookFailure
	PostHookFailure *hookFailure
}

func (s *step) kind() kind {
	return stepKind
}

type concept struct {
	CptStep *step
	Items   []item
}

func (c *concept) kind() kind {
	return conceptKind
}

type comment struct {
	Text string
}

func (c *comment) kind() kind {
	return commentKind
}

type result struct {
	Status        status
	StackTrace    string
	Screenshot    string
	ErrorMessage  string
	ExecTime      string
	SkippedReason string
	Messages      []string
}

type searchIndex struct {
	Tags  map[string][]string `json:"tags"`
	Specs map[string][]string `json:"specs"`
}

type status int

const (
	pass status = iota
	fail
	skip
	notExecuted
)

var parsedTemplates = make(map[string]*template.Template, 0)

// Any new templates that are added in file `templates.go` should be registered here
var templates = []string{bodyFooterTag, reportOverviewTag, sidebarDiv, congratsDiv, hookFailureDiv, tagsDiv, messageDiv, skippedReasonDiv,
	specsStartDiv, specsItemsContainerDiv, specsItemsContentsDiv, specHeaderStartTag, scenarioContainerStartDiv, scenarioHeaderStartDiv, specCommentsAndTableTag,
	htmlPageStartTag, headerEndTag, mainEndTag, endDiv, conceptStartDiv, stepStartDiv, stepMetaDiv, stepBodyDiv, stepFailureDiv, stepEndDiv, conceptSpan,
	contextOrTeardownStartDiv, commentSpan, conceptStepsStartDiv, nestedConceptDiv, htmlPageEndWithJS, specErrorDiv,
}

func init() {
	var encodeNewLine = func(s string) string {
		return strings.Replace(s, "\n", "<br/>", -1)
	}
	var parseMarkdown = func(args ...interface{}) string {
		s := blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...)))
		return string(s)
	}
	var sanitizeHTML = func(s string) string {
		var b bytes.Buffer
		var html = bluemonday.UGCPolicy().SanitizeBytes([]byte(s))
		b.Write(html)
		return b.String()
	}
	var funcs = template.FuncMap{"parseMarkdown": parseMarkdown, "sanitize": sanitizeHTML, "escapeHTML": template.HTMLEscapeString, "encodeNewLine": encodeNewLine}
	for _, tmpl := range templates {
		t, err := template.New("Reports").Funcs(funcs).Parse(tmpl)
		if err != nil {
			log.Fatalf(err.Error())
		}
		parsedTemplates[tmpl] = t
	}
}

func execTemplate(tmplName string, w io.Writer, data interface{}) {
	tmpl := parsedTemplates[tmplName]
	if tmpl == nil {
		log.Fatal(tmplName)
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

// ProjectRoot is root dir of current project
var ProjectRoot string

// GenerateReports generates HTML report in the given report dir location
func GenerateReports(suiteRes *gm.ProtoSuiteResult, reportDir string) error {
	f, err := os.Create(filepath.Join(reportDir, "index.html"))
	if err != nil {
		return err
	}
	if suiteRes.GetPreHookFailure() != nil {
		overview := toOverview(suiteRes, nil)
		generateOverview(overview, f)
		execTemplate(hookFailureDiv, f, toHookFailure(suiteRes.GetPreHookFailure(), "Before Suite"))
		if suiteRes.GetPostHookFailure() != nil {
			execTemplate(hookFailureDiv, f, toHookFailure(suiteRes.GetPostHookFailure(), "After Suite"))
		}
		generatePageFooter(overview, f)
	} else {
		var wg sync.WaitGroup
		wg.Add(1)
		go generateIndexPage(suiteRes, f, &wg)
		specRes := suiteRes.GetSpecResults()
		for _, res := range specRes {
			relPath, _ := filepath.Rel(ProjectRoot, res.GetProtoSpec().GetFileName())
			CreateDirectory(filepath.Join(reportDir, filepath.Dir(relPath)))
			sf, err := os.Create(filepath.Join(reportDir, toHTMLFileName(res.GetProtoSpec().GetFileName(), ProjectRoot)))
			if err != nil {
				return err
			}
			wg.Add(1)
			go generateSpecPage(suiteRes, res, sf, &wg)
		}
		wg.Wait()
	}
	err = generateSearchIndex(suiteRes, reportDir)
	if err != nil {
		return err
	}
	return nil
}

func newSearchIndex() *searchIndex {
	var i searchIndex
	i.Tags = make(map[string][]string)
	i.Specs = make(map[string][]string)
	return &i
}

func (i *searchIndex) hasValueForTag(tag string, spec string) bool {
	for _, s := range i.Tags[tag] {
		if s == spec {
			return true
		}
	}
	return false
}

func (i *searchIndex) hasSpec(specHeading string, specFileName string) bool {
	for _, s := range i.Specs[specHeading] {
		if s == specFileName {
			return true
		}
	}
	return false
}

func containsParseErrors(errors []error) bool {
	for _, e := range errors {
		if e.(buildError).isParseError() {
			return true
		}
	}
	return false
}

func generateSearchIndex(suiteRes *gm.ProtoSuiteResult, reportDir string) error {
	CreateDirectory(filepath.Join(reportDir, "js"))
	f, err := os.Create(filepath.Join(reportDir, "js", "search_index.js"))
	if err != nil {
		return err
	}
	index := newSearchIndex()
	for _, r := range suiteRes.GetSpecResults() {
		spec := r.GetProtoSpec()
		specFileName := toHTMLFileName(spec.GetFileName(), ProjectRoot)
		for _, t := range spec.GetTags() {
			if !index.hasValueForTag(t, specFileName) {
				index.Tags[t] = append(index.Tags[t], specFileName)
			}
		}
		var addTagsFromScenario = func(s *gm.ProtoScenario) {
			for _, t := range s.GetTags() {
				if !index.hasValueForTag(t, specFileName) {
					index.Tags[t] = append(index.Tags[t], specFileName)
				}
			}
		}
		for _, i := range spec.GetItems() {
			if s := i.GetScenario(); s != nil {
				addTagsFromScenario(s)
			}
			if tds := i.GetTableDrivenScenario(); tds != nil {
				addTagsFromScenario(tds.GetScenario())
			}
		}
		specHeading := spec.GetSpecHeading()
		if !index.hasSpec(specHeading, specFileName) {
			index.Specs[specHeading] = append(index.Specs[specHeading], specFileName)
		}
	}
	s, err := json.Marshal(index)
	if err != nil {
		return err
	}
	f.WriteString(fmt.Sprintf("var index = %s;", s))
	return nil
}

func generateIndexPage(suiteRes *gm.ProtoSuiteResult, w io.Writer, wg *sync.WaitGroup) {
	defer wg.Done()
	overview := toOverview(suiteRes, nil)
	generateOverview(overview, w)
	if suiteRes.GetPostHookFailure() != nil {
		execTemplate(hookFailureDiv, w, toHookFailure(suiteRes.GetPostHookFailure(), "After Suite"))
	}
	execTemplate(specsStartDiv, w, nil)
	execTemplate(sidebarDiv, w, toSidebar(suiteRes, nil))
	if !suiteRes.GetFailed() {
		execTemplate(congratsDiv, w, nil)
	}
	execTemplate(endDiv, w, nil)
	generatePageFooter(overview, w)
}

func generateSpecPage(suiteRes *gm.ProtoSuiteResult, specRes *gm.ProtoSpecResult, w io.Writer, wg *sync.WaitGroup) {
	defer wg.Done()
	overview := toOverview(suiteRes, specRes)

	generateOverview(overview, w)

	if suiteRes.GetPreHookFailure() != nil {
		execTemplate(hookFailureDiv, w, toHookFailure(suiteRes.GetPreHookFailure(), "Before Suite"))
	}

	if suiteRes.GetPostHookFailure() != nil {
		execTemplate(hookFailureDiv, w, toHookFailure(suiteRes.GetPostHookFailure(), "After Suite"))
	}

	if suiteRes.GetPreHookFailure() == nil {
		execTemplate(specsStartDiv, w, nil)
		execTemplate(sidebarDiv, w, toSidebar(suiteRes, specRes))
		generateSpecDiv(w, specRes)
		execTemplate(endDiv, w, nil)
	}
	generatePageFooter(overview, w)
}

func generateOverview(overview *overview, w io.Writer) {
	execTemplate(htmlPageStartTag, w, overview)
	execTemplate(reportOverviewTag, w, overview)
}

func generatePageFooter(overview *overview, w io.Writer) {
	execTemplate(endDiv, w, nil)
	execTemplate(mainEndTag, w, nil)
	execTemplate(bodyFooterTag, w, nil)
	execTemplate(htmlPageEndWithJS, w, overview)
}

func generateSpecDiv(w io.Writer, res *gm.ProtoSpecResult) {
	specHeader := toSpecHeader(res)
	spec := toSpec(res)

	execTemplate(specHeaderStartTag, w, specHeader)
	execTemplate(tagsDiv, w, specHeader)
	execTemplate(headerEndTag, w, nil)
	execTemplate(specsItemsContainerDiv, w, nil)
	if containsParseErrors(spec.Errors) {
		execTemplate(specErrorDiv, w, spec)
		execTemplate(endDiv, w, nil)
		return
	}

	if spec.BeforeHookFailure != nil {
		execTemplate(hookFailureDiv, w, spec.BeforeHookFailure)
	}

	execTemplate(specsItemsContentsDiv, w, nil)
	execTemplate(specCommentsAndTableTag, w, spec)

	if spec.BeforeHookFailure == nil {
		for _, scn := range spec.Scenarios {
			generateScenario(w, scn)
		}
	}

	execTemplate(endDiv, w, nil)
	execTemplate(endDiv, w, nil)

	if spec.AfterHookFailure != nil {
		execTemplate(hookFailureDiv, w, spec.AfterHookFailure)
	}

	execTemplate(endDiv, w, nil)
}

func generateScenario(w io.Writer, scn *scenario) {
	execTemplate(scenarioContainerStartDiv, w, scn)
	execTemplate(scenarioHeaderStartDiv, w, scn)
	execTemplate(tagsDiv, w, scn)
	execTemplate(endDiv, w, nil)
	if scn.BeforeHookFailure != nil {
		execTemplate(hookFailureDiv, w, scn.BeforeHookFailure)
	}

	generateItems(w, scn.Contexts, generateContextOrTeardown)
	generateItems(w, scn.Items, generateItem)
	generateItems(w, scn.Teardown, generateContextOrTeardown)

	if scn.AfterHookFailure != nil {
		execTemplate(hookFailureDiv, w, scn.AfterHookFailure)
	}
	execTemplate(endDiv, w, nil)
}

func generateItems(w io.Writer, items []item, predicate func(w io.Writer, item item)) {
	for _, item := range items {
		predicate(w, item)
	}
}

func generateContextOrTeardown(w io.Writer, item item) {
	execTemplate(contextOrTeardownStartDiv, w, nil)
	generateItem(w, item)
	execTemplate(endDiv, w, nil)
}

func generateItem(w io.Writer, item item) {
	switch item.kind() {
	case stepKind:
		execTemplate(stepStartDiv, w, item.(*step))
		execTemplate(stepBodyDiv, w, item.(*step))

		if item.(*step).PreHookFailure != nil {
			execTemplate(hookFailureDiv, w, item.(*step).PreHookFailure)
		}

		stepRes := item.(*step).Res
		if stepRes.Status == fail && stepRes.ErrorMessage != "" && stepRes.StackTrace != "" {
			execTemplate(stepFailureDiv, w, stepRes)
		}

		if item.(*step).PostHookFailure != nil {
			execTemplate(hookFailureDiv, w, item.(*step).PostHookFailure)
		}
		execTemplate(messageDiv, w, stepRes)
		execTemplate(stepEndDiv, w, item.(*step))
		if stepRes.Status == skip && stepRes.SkippedReason != "" {
			execTemplate(skippedReasonDiv, w, stepRes)
		}
	case commentKind:
		execTemplate(commentSpan, w, item.(*comment))
	case conceptKind:
		execTemplate(conceptStartDiv, w, item.(*concept).CptStep)
		execTemplate(conceptSpan, w, nil)
		execTemplate(stepBodyDiv, w, item.(*concept).CptStep)
		execTemplate(stepEndDiv, w, item.(*concept).CptStep)
		execTemplate(conceptStepsStartDiv, w, nil)
		generateItems(w, item.(*concept).Items, generateItem)
		execTemplate(endDiv, w, nil)
	}
}

// CreateDirectory creates given directory if it doesn't exist
func CreateDirectory(dir string) {
	if common.DirExists(dir) {
		return
	}
	if err := os.MkdirAll(dir, common.NewDirectoryPermissions); err != nil {
		fmt.Printf("Failed to create directory %s: %s\n", dir, err)
		os.Exit(1)
	}
}
