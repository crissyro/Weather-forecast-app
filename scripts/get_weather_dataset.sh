#!/bin/bash

## @file download_weather_data.sh
## @brief Скрипт для загрузки и обработки данных о погоде
## @details
## Скрипт выполняет следующие действия:
## 1. Скачивает архив с данными о погоде с Kaggle
## 2. Распаковывает архив
## 3. Удаляет исходный архив
## 4. Создает необходимые директории
##
## @author Moroz Roman
## @version 1.0
## @date 2025-03-25

set -euo pipefail 

DOWNLOAD_URL="https://www.kaggle.com/api/v1/datasets/download/nelgiriyewithana/global-weather-repository"
TARGET_DIR="${HOME}/Weather-forecast-app/data"
ARCHIVE_NAME="global-weather-repository.zip"
EXTRACT_DIR="${TARGET_DIR}"
DB_NAME="${TARGET_DIR}/state.db"

## @fn main
## @brief Основная функция скрипта
main() {
    echo "Starting weather data processing..."
    
    check_dependencies
    
    create_directories
    
    download_dataset
    
    extract_data
    
    cleanup
    
    echo "Processing completed successfully!"
    echo "Extracted data: ${EXTRACT_DIR}"
}

## @fn check_dependencies
## @brief Проверка необходимых зависимостей
check_dependencies() {
    local dependencies=("curl" "unzip")
    for cmd in "${dependencies[@]}"; do
        if ! command -v "${cmd}" &> /dev/null; then
            echo "Error: ${cmd} is not installed."
            exit 1
        fi
    done
}

## @fn create_directories
## @brief Создание целевых директорий
create_directories() {
    mkdir -p "${TARGET_DIR}"
    mkdir -p "${EXTRACT_DIR}"
}

## @fn download_dataset
## @brief Загрузка архива с данными
download_dataset() {
    echo "Downloading dataset..."
    curl -L -o "${TARGET_DIR}/${ARCHIVE_NAME}" "${DOWNLOAD_URL}" || {
        echo "Failed to download dataset"
        exit 1
    }
}

## @fn extract_data
## @brief Распаковка архива
extract_data() {
    echo "Extracting data..."
    unzip -q "${TARGET_DIR}/${ARCHIVE_NAME}" -d "${EXTRACT_DIR}" || {
        echo "Failed to extract archive"
        exit 1
    }
}

## @fn cleanup
## @brief Очистка временных файлов
cleanup() {
    echo "Cleaning up..."
    rm -f "${TARGET_DIR}/${ARCHIVE_NAME}"
}

trap 'echo "Script interrupted"; cleanup; exit 1' INT TERM

main
