# PRD: Delcom Blog Backend (Golang Gin)

## Product Overview
Delcom Blog adalah platform publikasi artikel multi-tier yang dibangun dengan arsitektur Clean Architecture dan Feature-First pada backend Golang Gin. Aplikasi ini menyediakan REST API untuk publikasi artikel, manajemen konten bagi penulis dan admin, sistem autentikasi berbasis JWT, serta AI Blog Generator yang terintegrasi dengan Search Engine dan OpenAI compatible API.

## Problem
Penulis membutuhkan platform publikasi artikel yang cepat, memiliki alur penulisan yang dibantu kecerdasan buatan dengan sitasi rujukan faktual dari mesin pencari, serta pemisahan hak akses berjenjang yang aman antara pengunjung (publik), penulis (author), dan administrator.

## Goal
- Menyediakan REST API backend yang tangguh, modular, dan teruji 100% statement coverage.
- Mengimplementasikan manajemen artikel (CRUD, slug unik otomatis, estimasi reading time, status publikasi instan).
- Mengintegrasikan AI Blog Generator dengan Search Engine rujukan web.
- Menyediakan dokumentasi API interaktif menggunakan Swagger UI.

## Target Users
| Persona | Kebutuhan |
|---|---|
| Publik (Guest) | Membaca artikel yang terbit, mencari artikel dengan kata kunci, memfilter kategori, melihat waktu baca. |
| Penulis (Author) | Mendaftar akun, mengelola draf dan artikel sendiri, memantau statistik penulisan, menghasilkan draf dengan AI Generator. |
| Administrator | Memantau seluruh akun pengguna dan seluruh arsip artikel di platform. |

## Core Features
| ID | Fitur | Prioritas |
|---|---|---|
| F-01 | Autentikasi & Profil (Register, Login JWT, Profile, Logout) | P0 |
| F-02 | Publikasi Publik (Kategori, Daftar Artikel Terbit, Detail Slug) | P0 |
| F-03 | Manajemen Artikel Penulis (CRUD, Slug Generator, Reading Time, Toggle Publish, Stats) | P0 |
| F-04 | Panel Administrasi (Daftar Pengguna, Arsip Seluruh Artikel, Artikel per Pengguna) | P1 |
| F-05 | AI Blog Generator (Search Engine Lookup + OpenAI Chat Completion + Citations) | P1 |
| F-06 | Dokumentasi Swagger UI interaktif & API Specification | P1 |

## Di Luar Scope
- Pengelolaan komentar pembaca (direncanakan untuk iterasi masa depan).
- Pembayaran atau langganan berbayar.
- Multi-bahasa selain Bahasa Indonesia / Inggris pada konten artikel.

## Success Metrics
- 100% statement test coverage pada seluruh paket backend `internal/...`.
- Seluruh endpoint merespons sesuai format JSON standar (`status`, `message`, `data`).
- Waktu respons API rata-rata < 200ms untuk endpoint non-AI.
