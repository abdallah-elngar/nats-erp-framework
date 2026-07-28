package template

import (
	"regexp"
	"strings"
)

// InheritanceManager يدير نظام الوراثة
type InheritanceManager struct {
	engine *Engine
}

// NewInheritanceManager ينشئ مدير وراثة جديد
func NewInheritanceManager(engine *Engine) *InheritanceManager {
	return &InheritanceManager{
		engine: engine,
	}
}

// ProcessTemplate يعالج القالب مع الوراثة
func (im *InheritanceManager) ProcessTemplate(content string) string {
	// معالجة extends
	if strings.Contains(content, "{{ extends ") {
		layoutName := im.extractLayoutName(content)
		if layoutName != "" {
			layoutContent, err := im.engine.loadLayout(layoutName)
			if err == nil {
				return im.inheritTemplate(content, layoutContent)
			}
		}
	}
	return content
}

// extractLayoutName يستخرج اسم التخطيط
func (im *InheritanceManager) extractLayoutName(content string) string {
	re := regexp.MustCompile(`{{\s*extends\s+"([^"]+)"\s*}}`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// inheritTemplate يدمج القالب مع التخطيط
func (im *InheritanceManager) inheritTemplate(content, layout string) string {
	// إزالة وسم extends من المحتوى
	content = regexp.MustCompile(`{{\s*extends\s+"[^"]+"\s*}}`).ReplaceAllString(content, "")
	
	// استخراج الكتل
	blocks := im.extractBlocks(content)
	
	// دمج الكتل مع التخطيط
	result := layout
	for name, blockContent := range blocks {
		// استبدال {{ block "name" }}...{{ endblock }} في التخطيط
		pattern := regexp.MustCompile(`{{\s*block\s+"` + regexp.QuoteMeta(name) + `"\s*}}[\s\S]*?{{\s*endblock\s*}}`)
		result = pattern.ReplaceAllString(result, blockContent)
	}
	
	return result
}

// extractBlocks يستخرج الكتل من المحتوى
func (im *InheritanceManager) extractBlocks(content string) map[string]string {
	blocks := make(map[string]string)
	
	re := regexp.MustCompile(`{{\s*block\s+"([^"]+)"\s*}}([\s\S]*?){{\s*endblock\s*}}`)
	matches := re.FindAllStringSubmatch(content, -1)
	
	for _, match := range matches {
		if len(match) > 2 {
			name := strings.TrimSpace(match[1])
			blockContent := strings.TrimSpace(match[2])
			blocks[name] = blockContent
		}
	}
	
	return blocks
}

// ExtractBlocksNames يستخرج أسماء الكتل
func (im *InheritanceManager) ExtractBlocksNames(content string) []string {
	var names []string
	re := regexp.MustCompile(`{{\s*block\s+"([^"]+)"\s*}}`)
	matches := re.FindAllStringSubmatch(content, -1)
	
	for _, match := range matches {
		if len(match) > 1 {
			names = append(names, match[1])
		}
	}
	
	return names
}

// HasBlock يتحقق من وجود كتلة
func (im *InheritanceManager) HasBlock(content, blockName string) bool {
	pattern := regexp.MustCompile(`{{\s*block\s+"` + regexp.QuoteMeta(blockName) + `"\s*}}`)
	return pattern.MatchString(content)
}

// GetBlockContent يحصل على محتوى الكتلة
func (im *InheritanceManager) GetBlockContent(content, blockName string) string {
	pattern := regexp.MustCompile(`{{\s*block\s+"` + regexp.QuoteMeta(blockName) + `"\s*}}([\s\S]*?){{\s*endblock\s*}}`)
	matches := pattern.FindStringSubmatch(content)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}