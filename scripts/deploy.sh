#!/bin/bash

# Production deployment script для Go приложения
# Использование: ./scripts/deploy.sh [up|down|migrate|logs|status]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
COMPOSE_FILE="$PROJECT_ROOT/docker-compose.prod.yml"

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Функции логирования
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Проверка наличия docker-compose
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker не установлен"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        log_error "docker-compose не установлен"
        exit 1
    fi
    
    log_success "Docker и docker-compose найдены"
}

# Сборка и запуск контейнеров
deploy_up() {
    log_info "Сборка production образов..."
    docker-compose -f "$COMPOSE_FILE" build
    
    log_info "Запуск production контейнеров..."
    docker-compose -f "$COMPOSE_FILE" up -d
    
    log_info "Ожидание готовности БД..."
    sleep 5
    
    log_info "Запуск миграций..."
    docker-compose -f "$COMPOSE_FILE" exec -T app /app/migrate up
    
    log_success "Production deployment завершен успешно!"
    log_info "API доступен на http://localhost:8080"
}

# Остановка контейнеров
deploy_down() {
    log_warning "Остановка production контейнеров..."
    docker-compose -f "$COMPOSE_FILE" down --remove-orphans
    log_success "Контейнеры остановлены"
}

# Запуск миграций
deploy_migrate() {
    log_info "Запуск миграций..."
    docker-compose -f "$COMPOSE_FILE" exec -T app /app/migrate up
    log_success "Миграции выполнены успешно"
}

# Откат миграций
deploy_migrate_down() {
    log_warning "Откат миграций..."
    docker-compose -f "$COMPOSE_FILE" exec -T app /app/migrate down
    log_success "Миграции откачены"
}

# Проверка статуса миграций
deploy_status() {
    log_info "Статус миграций:"
    docker-compose -f "$COMPOSE_FILE" exec -T app /app/migrate status
}

# Просмотр логов
deploy_logs() {
    log_info "Логи API (Ctrl+C для выхода):"
    docker-compose -f "$COMPOSE_FILE" logs -f app
}

# Просмотр логов БД
deploy_logs_db() {
    log_info "Логи БД (Ctrl+C для выхода):"
    docker-compose -f "$COMPOSE_FILE" logs -f db
}

# Проверка здоровья приложения
deploy_health() {
    log_info "Проверка здоровья приложения..."
    
    if docker-compose -f "$COMPOSE_FILE" ps app | grep -q "Up"; then
        log_success "API контейнер запущен"
    else
        log_error "API контейнер не запущен"
        return 1
    fi
    
    if docker-compose -f "$COMPOSE_FILE" ps db | grep -q "Up"; then
        log_success "БД контейнер запущен"
    else
        log_error "БД контейнер не запущен"
        return 1
    fi
    
    # Проверка доступности API
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        log_success "API доступен"
    else
        log_warning "API недоступен (может быть в процессе запуска)"
    fi
}

# Справка
show_help() {
    cat << EOF
Production Deployment Script для Trainer API

Использование: $0 [COMMAND]

Команды:
    up              - Сборка и запуск production контейнеров с миграциями
    down            - Остановка production контейнеров
    migrate         - Запуск миграций БД
    migrate-down    - Откат миграций БД
    status          - Проверка статуса миграций
    logs            - Просмотр логов API
    logs-db         - Просмотр логов БД
    health          - Проверка здоровья приложения
    help            - Показать эту справку

Примеры:
    $0 up           # Полный деплой
    $0 logs         # Просмотр логов
    $0 migrate      # Запуск миграций
    $0 down         # Остановка

EOF
}

# Основная логика
main() {
    check_docker
    
    case "${1:-help}" in
        up)
            deploy_up
            ;;
        down)
            deploy_down
            ;;
        migrate)
            deploy_migrate
            ;;
        migrate-down)
            deploy_migrate_down
            ;;
        status)
            deploy_status
            ;;
        logs)
            deploy_logs
            ;;
        logs-db)
            deploy_logs_db
            ;;
        health)
            deploy_health
            ;;
        help)
            show_help
            ;;
        *)
            log_error "Неизвестная команда: $1"
            show_help
            exit 1
            ;;
    esac
}

main "$@"
