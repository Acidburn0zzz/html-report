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

const bodyHeaderTag = `<header class="top">
  <div class="header">
    <div class="container">
      <div class="logo">
        <a href="{{.BasePath}}"><img src="{{.BasePath}}images/logo.png" alt="Report logo"></a>
      </div>
      <h2 class="project">Project: {{.ProjectName}}</h2>
    </div>
  </div>
</header>`

const bodyFooterTag = `<footer class="footer">
  <div class="container">
    <p>Generated by Gauge HTML Report</p>
  </div>
</footer>`

const reportOverviewTag = `<div class="report-overview">
  <div class="report_chart">
    <div class="chart">
      <svg></svg>
    </div>
    <div class="total-specs"><span class="value">{{.TotalSpecs}}</span><span class="txt">Total specs</span></div>
  </div>
  <div class="report_test-results">
    <ul>
      <li class="fail spec-filter" data-status="failed"><span class="value">{{.Failed}}</span><span class="txt">Failed</span></li>
      <li class="pass spec-filter" data-status="passed"><span class="value">{{.Passed}}</span><span class="txt">Passed</span></li>
      <li class="skip spec-filter" data-status="skipped"><span class="value">{{.Skipped}}</span><span class="txt">Skipped</span></li>
    </ul>
  </div>
  <div class="report_details">
    <ul>
      <li>
        <label>Environment </label>
        <span>
          <span>{{.Env}}</span>
          <a href="#" id="view-env">View Env Vars</a>
          <div class="env-vars hidden">
            <div class="popup-header">
              <a style="float: right;" class="close-popup">
                <i class="fa fa-times" aria-hidden="true"></i>
              </a>
            </div>
            <table>
              <thead>
                <tr>
                  <th>KEY=Value</th>
                </tr>
              </thead>
              <tbody>
                {{range .EnvVars}}<tr><td>{{.}}</td></tr>{{end}}
              </tbody>
            </table>
          </div>
        </span>
      </li>
      {{if .Tags}}
      <li>
        <label>Tags </label>
        <span>{{.Tags}}</span>
      </li>
      {{end}}
      <li>
        <label>Success Rate </label>
        <span>{{.SuccRate}}%</span>
      </li>
      <li>
        <label>Total Time </label>
        <span>{{.ExecTime}}</span>
      </li>
      <li>
        <label>Generated On </label>
        <span>{{.Timestamp}}</span>
      </li>
    </ul>
  </div>
</div>`

//TODO: 1. Filtering based on search query
const sidebarDiv = `{{if not .IsBeforeHookFailure}}<aside class="sidebar">
  <h3 class="title">Specifications</h3>

  <div class="searchbar">
    <input id="searchSpecifications" placeholder="Type specification or tag name" type="text" />
    <i class="fa fa-search"></i>
  </div>

  <div id="listOfSpecifications">
    <ul id="scenarios" class="spec-list">
    {{range $index, $specMeta := .Specs}}
      <a href="{{.ReportFile}}">
        {{if $specMeta.Failed}} <li class='failed spec-name'>
        {{else if $specMeta.Skipped}} <li class='skipped spec-name'>
        {{else}} <li class='passed spec-name'>
        {{end}}
          <span id="scenarioName" class="scenarioname">{{$specMeta.SpecName}}</span>
          <span id="time" class="time">{{$specMeta.ExecTime}}</span>
        </li>
      </a>
      {{end}}
    </ul>
  </div>
</aside>{{end}}`

//TODO: Hide if pre/post hook failed
const congratsDiv = `
  <div class="congratulations details">
    <p>Congratulations! You've gone all <span class="green">green</span> and saved the environment!</p>
  </div>`

//TODO 1. Change text on toggle collapse
//     2. Check for collapsible
const hookFailureDiv = `<div class="error-container failed">
  <div collapsable class="error-heading">{{.HookName}} Failed:<span class="error-message"> {{.ErrMsg}}</span></div>
  <div class="toggleShow" data-toggle="collapse" data-target="#hookFailureDetails">
    <span>[Show details]</span>
  </div>
  <div class="exception-container" id="hookFailureDetails">
      <div class="exception">
        <pre class="stacktrace">{{.StackTrace}}</pre>
      </div>
      {{if .Screenshot}}<div class="screenshot-container">
        <a href="data:image/png;base64,{{.Screenshot}}" rel="lightbox">
          <img src="data:image/png;base64,{{.Screenshot}}" class="screenshot-thumbnail" />
        </a>
      </div>{{end}}
  </div>
</div>`

const tagsDiv = `{{if .Tags}}<div class="tags scenario_tags contentSection">
  <strong>Tags:</strong>
  {{range .Tags}}<span> {{.}}</span>{{end}}
</div>{{end}}`

//TODO 1. Format message to convert newlines to <br>
const messageDiv = `{{if .Messages}}<div class="message-container">
  {{range .Messages}}<p class="step-message">{{.}}</p>{{end}}
</div>{{end}}`

const skippedReasonDiv = `<div class="message-container">
  <h4 class="skipReason">Skipped Reason: {{.SkippedReason}}</h4>
</div>`

const specsStartDiv = `<div class="specifications">`

//TODO: Hide this if there is a pre hook failure
const specContainerStartDiv = `<div id="specificationContainer" class="details">`

const specsItemsContainerDiv = `<div id="specItemsContainer">`
const specsItemsContentsDiv = `<div class="content">`

const specHeaderStartTag = `<header class="curr-spec">
  <h3 class="spec-head" title="{{.FileName}}">{{.SpecName}}</h3>
  <div class="spec-filename">
    <label for="specFileName">File Path</label>
    <input id="specFileName" value="{{.FileName}}" readonly>
    <button class="clipboard-btn" data-clipboard-target="#specFileName" title="Copy to Clipboard">
        <i class="fa fa-clipboard" aria-hidden="true" title="Copy to Clipboard"></i>
    </button>
  </div>
  <span class="time">{{.ExecTime}}</span>`

const scenarioContainerStartDiv = `{{if eq .ExecStatus 0}}<div class='scenario-container passed{{if gt .TableRowIndex 0}} hidden{{end}}'{{if ne .TableRowIndex -1}} data-tablerow={{.TableRowIndex}}{{end}}>
{{else if eq .ExecStatus 1}}<div class='scenario-container failed{{if gt .TableRowIndex 0}} hidden{{end}}'{{if ne .TableRowIndex -1}}  data-tablerow={{.TableRowIndex}}{{end}}>
{{else}}<div class='scenario-container skipped{{if gt .TableRowIndex 0}} hidden{{end}}'{{if ne .TableRowIndex -1}}  data-tablerow={{.TableRowIndex}}{{end}}>{{end}}`

const scenarioHeaderStartDiv = `<div class="scenario-head">
  <h3 class="head borderBottom">{{.Heading}}</h3>
  <span class="time">{{.ExecTime}}</span>`

// TODO 1. Implement onclick of row
//      2. Set class as 'selected' on click
//      3. Convert comments to markdown. Check that this adds <p> tag
//         for every new line which is not happening right now.
const specCommentsAndTableTag = `{{range .CommentsBeforeTable}}<span>{{.}}</span>{{end}}
{{if .Table}}<table class="data-table">
  <tr>
    {{range .Table.Headers}}<th>{{.}}</th>{{end}}
  </tr>
  <tbody data-rowCount={{len .Table.Rows}}>
    {{range $index, $row := .Table.Rows}}
      {{if eq $row.Res 0}}<tr class='row-selector passed{{if eq $index 0}} selected{{end}}' data-rowIndex={{$index}}>
      {{else if eq $row.Res 1}}<tr class='row-selector failed{{if eq $index 0}} selected{{end}}' data-rowIndex={{$index}}>
      {{else}}<tr class='row-selector skipped{{if eq $index 0}} selected{{end}}' data-rowIndex={{$index}}>
      {{end}}
        {{range $row.Cells}}<td>{{.}}</td>{{end}}
    </tr>
    {{end}}
  </tbody>
</table>{{end}}
{{range .CommentsAfterTable}}<span>{{.}}</span>{{end}}`

// Common HTML tags templates
const htmlStartTag = `<!doctype html>
<html>`

const htmlEndTag = `</html>`

//TODO: Move JS includes at the end of body
const pageHeaderTag = `<head>
  <meta http-equiv="X-UA-Compatible" content="IE=9; IE=8; IE=7; IE=EDGE" />
  <meta charset="utf-8"/>
  <title>Gauge Test Results</title>
  <link rel="shortcut icon" type="image/x-icon" href="{{.BasePath}}images/favicon.ico">
  <link rel="stylesheet" type="text/css" href="{{.BasePath}}css/open-sans.css">
  <link rel="stylesheet" type="text/css" href="{{.BasePath}}css/font-awesome.css">
  <link rel="stylesheet" type="text/css" href="{{.BasePath}}css/normalize.css" />
  <link rel="stylesheet" type="text/css" href="{{.BasePath}}css/style.css" />
  <script src="{{.BasePath}}js/d3.min.js" charset="utf-8"></script>
  <script src="{{.BasePath}}js/nv.d3.min.js" charset="utf-8"></script>
  <script src="{{.BasePath}}js/lightbox.js"></script>
  <script src="{{.BasePath}}js/chart.js" type="text/javascript"></script>
  <script src="{{.BasePath}}js/jquery-3.1.0.min.js" type="text/javascript"></script>
  <script src="{{.BasePath}}js/auto-complete.min.js" type="text/javascript"></script>
  <script src="{{.BasePath}}js/clipboard.min.js" type="text/javascript"></script>
  <script src="{{.BasePath}}js/search_index.js" type="text/javascript"></script>
  <script src="{{.BasePath}}js/main.js" type="text/javascript"></script>
  <script type="text/javascript">
    createChart({{.Passed}},{{.Failed}},{{.Skipped}});
    var loadingImage = '{{.BasePath}}images/loading.gif';
    var closeButton = '{{.BasePath}}images/close.gif';
  </script>
</head>`

const headerEndTag = `</header>`

const bodyStartTag = `<body>`

const bodyEndTag = `</body>`

const mainStartTag = `<main class="main-container">`

const mainEndTag = `</main>`

const containerStartDiv = `<div class="container">`

const endDiv = `</div>`

const conceptStartDiv = `<div class='step concept'>` + stepMetaDiv
const stepStartDiv = `<div class='step'>` + stepMetaDiv

const stepMetaDiv = `
  {{if ne .Res.Status 2}}
  <h5 class='execution-time'>
  <span class='time'>Execution Time : {{.Res.ExecTime}}</span>
  </h5>
  {{end}}
    {{if eq .Res.Status 0}}<div class='step-info passed'>
    {{else if eq .Res.Status 1}}<div class='step-info failed'>
    {{else if eq .Res.Status 2}}<div class='step-info skipped'>
    {{else if eq .Res.Status 3}}<div class='step-info not-executed'>
    {{end}}
    <ul>
      <li class='step'>
        <div class='step-txt'>`

//TODO:
//  2. Print Pre/Post hook failures, Step failure
const stepBodyDiv = `
{{define "Table"}}<table>
  <tr>
    {{range .Table.Headers}}<th>{{.}}</th>{{end}}
  </tr>
  <tbody>
    {{range .Table.Rows}}
    <tr>{{range .Cells}}<td>{{.}}</td>{{end}}</tr>
    {{end}}
  </tbody>
</table>
{{end}}
{{range .Fragments}}
  {{if eq .FragmentKind 0}}
    <span>
      {{.Text}}
    </span>
  {{else if eq .FragmentKind 1 2}}
    <span class='parameter'>"{{.Text}}"</span>
  {{else if eq .FragmentKind 3}}
    <span class="hoverable">&lt;{{.Name}}&gt;</span>
    <div class="hovercard">{{.Text}}</div>
  {{else if eq .FragmentKind 4}}
    <span class="hoverable">&lt;{{.Name}}&gt;</span>
    <div class="hovercard">{{template "Table" .}}</div>
  {{else if eq .FragmentKind 5}}
    <div class='inline-table'>
      <div>
        {{template "Table" .}}
      </div>
    </div>
  {{end}}
{{end}}
</div>`

const stepFailureDiv = `<div class="error-container failed">
  <div class="exception-container">
      <div class="exception">
        <h4 class="error-message">
          <pre>{{.ErrorMessage}}</pre>
        </h4>
        <pre class="stacktrace">{{.StackTrace}}</pre>
      </div>
      {{if .Screenshot}}<div class="screenshot-container">
        <a href="data:image/png;base64,{{.Screenshot}}" rel="lightbox">
          <img src="data:image/png;base64,{{.Screenshot}}" class="screenshot-thumbnail" />
        </a>
      </div>{{end}}
  </div>
</div>`

const stepEndDiv = `</li></ul></div></div>`

const conceptSpan = `<i class="fa fa-plus-square" aria-hidden="true"></i>`

const contextOrTeardownStartDiv = `<div class='context-step'>`

//TODO: 1. Show comments in markdown style
const commentSpan = `<span>{{.Text}}</span>`

const conceptStepsStartDiv = `<div class='concept-steps'>`

const nestedConceptDiv = `<div class="nested concept-steps">`
