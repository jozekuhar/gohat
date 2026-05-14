set dotenv-load

# Those regex are used to sort tailwind classes but they can also map all 
classRegex := 'Class\(\s*"([^\"]+)"\s*,?\s*\)'

# @format: The "Code-Quality Pipeline"
# 1. SORT:  Use rustywind with custom regex to organize Tailwind classes in Go Class(...) & structs.
# 2. CLEAN: Run gofmt to fix any whitespace/spacing messes the regex left behind.
# 3. JUDGE: Run golangci-lint to ensure the shuffle didn't break Go syntax or logic.
@format:
    pnpm rustywind --write --custom-regex 'Class\(\s*"([^\"]+)"\s*,?\s*\)' ./internal/view/
    gofmt -w .
    golangci-lint run

@run:
    docker-compose up -d
    pnpm outdated --json
    air \
        --build.delay "100" \
        --build.exclude_dir "./static/" \
        --build.include_ext "go,ts,.env" \
        --build.cmd " \
            pnpm vite build; \
            just format; \
            go build -o ./tmp/main ./cmd/main/main.go; \
        "

@snippet:
    go run ./cmd/snippet/main.go

@goose command optional1="" optional2="":
    GOOSE_DRIVER="postgres" GOOSE_DBSTRING="$DATABASE_URL" GOOSE_MIGRATION_DIR="./migrations/" goose {{command}} {{optional1}} {{optional2}}
