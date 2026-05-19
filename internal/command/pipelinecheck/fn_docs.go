// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2026-present Guance, Inc.

package pipelinecheck

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/GuanceCloud/pipeline-go/ptinput/funcs"
)

const defaultSnippetLen = 260

type functionDocResult struct {
	OK        bool          `json:"ok"`
	Mode      string        `json:"mode"`
	Query     string        `json:"query,omitempty"`
	Language  string        `json:"language"`
	Count     int           `json:"count"`
	Functions []functionDoc `json:"functions,omitempty"`
	Errors    []string      `json:"errors,omitempty"`
}

type functionDoc struct {
	Name        string `json:"name"`
	Language    string `json:"language"`
	File        string `json:"file"`
	Signature   string `json:"signature,omitempty"`
	Description string `json:"description,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Snippet     string `json:"snippet,omitempty"`
	Markdown    string `json:"markdown,omitempty"`

	searchScore int
}

func executeFunctionDocs(cfg config) (functionDocResult, int) {
	res := functionDocResult{
		Language: cfg.functionLang,
	}

	docs, err := loadFunctionDocs(cfg.functionLang)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		return res, 1
	}

	switch {
	case cfg.functionDoc != "":
		res.Mode = "doc"
		res.Query = cfg.functionDoc
		res.Functions = findFunctionDocs(docs, cfg.functionDoc)
		for i := range res.Functions {
			res.Functions[i].Markdown = docsMarkdown(docs, res.Functions[i])
		}
		if len(res.Functions) == 0 {
			res.Errors = append(res.Errors, fmt.Sprintf("function doc %q not found", cfg.functionDoc))
			return res, 1
		}
	case cfg.searchFunctions != "":
		res.Mode = "search"
		res.Query = cfg.searchFunctions
		res.Functions = searchFunctionDocs(docs, cfg.searchFunctions, cfg.functionLimit)
	case cfg.listFunctions:
		res.Mode = "list"
		res.Functions = compactFunctionDocs(limitFunctionDocs(docs, cfg.functionLimit))
	default:
		res.Errors = append(res.Errors, "no function doc mode selected")
		return res, 1
	}

	for i := range res.Functions {
		res.Functions[i].searchScore = 0
	}
	res.Count = len(res.Functions)
	res.OK = true
	return res, 0
}

func loadFunctionDocs(lang string) ([]functionDoc, error) {
	origLang := lang
	lang = normalizeDocLang(lang)
	if lang == "" {
		return nil, fmt.Errorf("invalid function language %q; use zh, en, or all", origLang)
	}

	var docs []functionDoc
	if lang == "zh" || lang == "all" {
		docs = append(docs, functionDocsFromPLDocs(funcs.PipelineFunctionDocs, "zh")...)
	}
	if lang == "en" || lang == "all" {
		docs = append(docs, functionDocsFromPLDocs(funcs.PipelineFunctionDocsEN, "en")...)
	}

	sort.Slice(docs, func(i, j int) bool {
		if docs[i].Name == docs[j].Name {
			return docs[i].Language < docs[j].Language
		}
		return docs[i].Name < docs[j].Name
	})

	return docs, nil
}

func functionDocsFromPLDocs(src map[string]*funcs.PLDoc, lang string) []functionDoc {
	docs := make([]functionDoc, 0, len(src))
	for name, plDoc := range src {
		if plDoc == nil {
			continue
		}
		docs = append(docs, functionDocFromPLDoc(name, lang, plDoc))
	}
	return docs
}

func functionDocFromPLDoc(name, lang string, plDoc *funcs.PLDoc) functionDoc {
	docName := normalizeFunctionDocName(name)
	doc := parseFunctionDoc(
		fmt.Sprintf("ptinput/funcs.PipelineFunctionDocs[%s:%s]", lang, name),
		docName,
		lang,
		plDoc.Doc,
	)
	if plDoc.Prototype != "" {
		doc.Signature = cleanDocValue(plDoc.Prototype)
	}
	if plDoc.Description != "" {
		doc.Description = cleanDocValue(plDoc.Description)
		doc.Summary = doc.Description
	}
	return doc
}

func normalizeFunctionDocName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.TrimSuffix(name, "()")
	return name
}

func normalizeDocLang(lang string) string {
	switch strings.ToLower(strings.TrimSpace(lang)) {
	case "", "zh", "cn", "zh-cn":
		return "zh"
	case "en", "eng", "en-us":
		return "en"
	case "all", "*":
		return "all"
	default:
		return ""
	}
}

func parseFunctionDoc(file, name, lang, markdown string) functionDoc {
	markdown = strings.ReplaceAll(markdown, "\r\n", "\n")
	doc := functionDoc{
		Name:     name,
		Language: lang,
		File:     file,
		Markdown: markdown,
	}

	for _, line := range strings.Split(markdown, "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "函数原型："):
			doc.Signature = cleanDocValue(strings.TrimPrefix(trimmed, "函数原型："))
		case strings.HasPrefix(trimmed, "Function prototype:"):
			doc.Signature = cleanDocValue(strings.TrimPrefix(trimmed, "Function prototype:"))
		case strings.HasPrefix(trimmed, "函数说明："):
			doc.Description = cleanDocValue(strings.TrimPrefix(trimmed, "函数说明："))
		case strings.HasPrefix(trimmed, "Function description:"):
			doc.Description = cleanDocValue(strings.TrimPrefix(trimmed, "Function description:"))
		}
	}

	doc.Summary = doc.Description
	if doc.Summary == "" {
		doc.Summary = firstUsefulLine(markdown)
	}

	return doc
}

func cleanDocValue(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "`")
	return strings.TrimSpace(s)
}

func firstUsefulLine(markdown string) string {
	inCode := false
	for _, line := range strings.Split(markdown, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inCode = !inCode
			continue
		}
		if inCode || trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "-") {
			continue
		}
		return cleanDocValue(trimmed)
	}
	return ""
}

func findFunctionDocs(docs []functionDoc, name string) []functionDoc {
	name = strings.ToLower(strings.TrimSpace(name))
	var ret []functionDoc
	for _, doc := range docs {
		if strings.ToLower(doc.Name) == name {
			ret = append(ret, doc)
		}
	}
	return ret
}

func searchFunctionDocs(docs []functionDoc, query string, limit int) []functionDoc {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return nil
	}

	var ret []functionDoc
	for _, doc := range docs {
		score := scoreFunctionDoc(doc, query)
		if score == 0 {
			continue
		}
		doc.searchScore = score
		doc.Snippet = makeSnippet(doc.Markdown, query, defaultSnippetLen)
		doc.Markdown = ""
		ret = append(ret, doc)
	}

	sort.Slice(ret, func(i, j int) bool {
		if ret[i].searchScore == ret[j].searchScore {
			if ret[i].Name == ret[j].Name {
				return ret[i].Language < ret[j].Language
			}
			return ret[i].Name < ret[j].Name
		}
		return ret[i].searchScore > ret[j].searchScore
	})

	return limitFunctionDocs(ret, limit)
}

func scoreFunctionDoc(doc functionDoc, query string) int {
	name := strings.ToLower(doc.Name)
	signature := strings.ToLower(doc.Signature)
	description := strings.ToLower(doc.Description)
	markdown := strings.ToLower(doc.Markdown)

	score := 0
	if name == query {
		score += 100
	}
	if strings.Contains(name, query) {
		score += 40
	}
	if strings.Contains(signature, query) {
		score += 25
	}
	if strings.Contains(description, query) {
		score += 20
	}
	if strings.Contains(markdown, query) {
		score += 5
	}
	return score
}

func limitFunctionDocs(docs []functionDoc, limit int) []functionDoc {
	if limit <= 0 || limit >= len(docs) {
		return docs
	}
	return docs[:limit]
}

func compactFunctionDocs(docs []functionDoc) []functionDoc {
	ret := make([]functionDoc, len(docs))
	copy(ret, docs)
	for i := range ret {
		ret[i].Markdown = ""
		ret[i].Snippet = ""
	}
	return ret
}

func makeSnippet(markdown, query string, maxLen int) string {
	plain := collapseSpace(stripMarkdownFences(markdown))
	lower := strings.ToLower(plain)
	idx := strings.Index(lower, strings.ToLower(query))
	if idx < 0 {
		if len(plain) <= maxLen {
			return plain
		}
		return strings.TrimSpace(plain[:maxLen]) + "..."
	}

	start := idx - maxLen/3
	if start < 0 {
		start = 0
	}
	end := start + maxLen
	if end > len(plain) {
		end = len(plain)
		start = end - maxLen
		if start < 0 {
			start = 0
		}
	}

	snippet := strings.TrimSpace(plain[start:end])
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(plain) {
		snippet += "..."
	}
	return snippet
}

func stripMarkdownFences(markdown string) string {
	var lines []string
	inCode := false
	for _, line := range strings.Split(markdown, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inCode = !inCode
			continue
		}
		if inCode {
			lines = append(lines, line)
			continue
		}
		lines = append(lines, strings.TrimPrefix(line, "### "))
	}
	return strings.Join(lines, " ")
}

func collapseSpace(s string) string {
	return strings.Join(strings.FieldsFunc(s, unicode.IsSpace), " ")
}

func docsMarkdown(docs []functionDoc, doc functionDoc) string {
	for _, candidate := range docs {
		if candidate.Name == doc.Name && candidate.Language == doc.Language {
			return candidate.Markdown
		}
	}
	return doc.Markdown
}
