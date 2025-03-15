package generator

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/pycnick/cagenerator/internal/config"
	"github.com/pycnick/cagenerator/internal/types"
	"github.com/pycnick/cagenerator/internal/utils"
)

// Generator handles project generation
type Generator struct {
	config      *config.Config
	moduleName  string
	templates   map[string][]string
	projectPath string
}

// New creates a new project generator
func New(cfg *config.Config, projectPath string) (*Generator, error) {
	// Проверяем наличие go.mod в указанном пути
	if _, err := os.Stat(filepath.Join(projectPath, "go.mod")); os.IsNotExist(err) {
		return nil, fmt.Errorf("go.mod not found in project directory: %s", projectPath)
	}

	// Переходим в директорию проекта
	if err := os.Chdir(projectPath); err != nil {
		return nil, fmt.Errorf("changing to project directory: %w", err)
	}

	moduleName, err := utils.GetModuleName()
	if err != nil {
		return nil, fmt.Errorf("getting module name: %w", err)
	}

	return &Generator{
		config:      cfg,
		moduleName:  moduleName,
		templates:   defaultTemplates(),
		projectPath: projectPath,
	}, nil
}

// Generate generates all project files
func (g *Generator) Generate() error {
	// Generate files for each entity
	for _, entity := range g.config.Entities {
		if err := g.generateEntityFiles(entity); err != nil {
			return fmt.Errorf("generating files for entity %s: %w", entity.Name, err)
		}
	}

	if err := utils.FormatProject(); err != nil {
		return fmt.Errorf("formatting project: %w", err)
	}

	return nil
}

// generateEntityFiles generates all files for a single entity
func (g *Generator) generateEntityFiles(entity types.Entity) error {
	// Преобразуем имена полей в snake_case
	var databaseFields []string
	for _, field := range entity.Fields {
		databaseFields = append(databaseFields, utils.CamelToSnake(field.Name))
	}

	// Преобразуем типы из YAML в Go и PostgreSQL
	for i, field := range entity.Fields {
		typeInfo, ok := types.SupportedTypes[field.Type]
		if !ok {
			return fmt.Errorf("unsupported entity type: %s", field.Type)
		}

		entity.Fields[i].GoType = typeInfo.GoType
		entity.Fields[i].PostgresType = typeInfo.PostgresType
		entity.Fields[i].DefaultGoValue = typeInfo.DefaultGoValue
	}

	data := map[string]interface{}{
		"ModuleName":      g.moduleName,
		"Entity":          entity,
		"LowerEntityName": strings.ToLower(entity.Name),
		"CommonFields":    entity.CommonFields(),
		"PrimaryField":    entity.PrimaryField(),
		"OwnerField":      entity.OwnerField(),
		"DatabaseFields":  databaseFields,
	}

	for layer, templateFiles := range g.templates {
		entityDir := filepath.Join(g.projectPath, "internal", strings.ToLower(entity.Name), layer)
		if err := utils.EnsureDirectory(entityDir); err != nil {
			return fmt.Errorf("creating directory for %s: %w", layer, err)
		}

		for _, templateFile := range templateFiles {
			if err := g.generateLayerFiles(layer, templateFile, data); err != nil {
				return fmt.Errorf("generating %s layer: %w", layer, err)
			}
		}
	}

	return nil
}

// Additional methods and implementations...

// defaultTemplates возвращает карту шаблонов по умолчанию
func defaultTemplates() map[string][]string {
	return map[string][]string{
		"domain": {
			"templates/domain.tmpl",
		},
		"repository": {
			"templates/repository/repository.tmpl",
			"templates/repository/tests.tmpl",
		},
		"usecase": {
			"templates/usecase/usecase.tmpl",
			"templates/usecase/tests.tmpl",
		},
		"delivery": {
			"templates/delivery/delivery.tmpl",
			"templates/delivery/tests.tmpl",
		},
		"migrations": {
			"templates/migration.tmpl",
		},
	}
}

// generateLayerFiles генерирует файлы для конкретного слоя
func (g *Generator) generateLayerFiles(layer, templateFile string, data map[string]interface{}) error {
	// Читаем содержимое шаблона
	tmplContent, err := utils.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("reading template file %s: %w", templateFile, err)
	}

	// Создаем и парсим шаблон
	funcMap := template.FuncMap{
		"lower": utils.Lower,
		"snake": utils.CamelToSnake,
	}

	tmpl, err := template.New(layer).Funcs(funcMap).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("parsing template %s: %w", templateFile, err)
	}

	// Особая обработка для миграций
	if layer == "migrations" {
		return g.generateMigrationFile(tmpl, data)
	}

	// Определяем имя выходного файла
	fileName := layer
	if strings.HasSuffix(templateFile, "tests.tmpl") {
		fileName = layer + "_test"
	}

	// Создаем выходной файл с учетом пути проекта
	entityDir := filepath.Join(g.projectPath, "internal", strings.ToLower(data["Entity"].(types.Entity).Name), layer)
	filePath := filepath.Join(entityDir, fileName+".go")

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", filePath, err)
	}
	defer file.Close()

	// Выполняем шаблон
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("executing template %s: %w", templateFile, err)
	}

	return nil
}

// generateMigrationFile генерирует SQL файл миграции
func (g *Generator) generateMigrationFile(tmpl *template.Template, data map[string]interface{}) error {
	// Создаем директорию migrations если её нет
	migrationsDir := filepath.Join(g.projectPath, "migrations")
	if err := utils.EnsureDirectory(migrationsDir); err != nil {
		return fmt.Errorf("creating migrations directory: %w", err)
	}

	// Формируем имя файла миграции
	entity := data["Entity"].(types.Entity)
	fileName := fmt.Sprintf("create_%s_table.sql", utils.CamelToSnake(entity.Name))
	filePath := filepath.Join(migrationsDir, fileName)

	// Создаем файл миграции
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("creating migration file %s: %w", filePath, err)
	}
	defer file.Close()

	// Выполняем шаблон
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("executing migration template: %w", err)
	}

	return nil
}
