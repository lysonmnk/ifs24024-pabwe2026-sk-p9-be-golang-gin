# Architecture
## System Overview
Browser -> Laravel (Blade + Livewire) -> Service/Action -> Repository -> MySQL
## Tech Stack (TERKUNCI)
Laravel (PHP 8.3+), Blade + Livewire + Tailwind CSS, MySQL, Pest/PHPUnit, Laravel Pint, Jenkins.
Menambah library = izin user + catat di docs/DECISIONS.md.
## Struktur Folder
app/Features/<Fitur>/{Actions,Services,Repositories,Policies,DTOs,Livewire}
app/Models, app/Http (Controllers, Requests)
resources/views, database/, tests/, contracts/, docs/
## Data Flow
Request -> FormRequest -> Controller -> Service/Action -> Repository -> DB -> Response sesuai kontrak
