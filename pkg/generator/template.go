package generator

const HTMLTemplate = `<!DOCTYPE html>
<html lang="en" data-theme="{{ .Config.Theme }}">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>{{ .Config.Title }}</title>
  <style>
    :root {
      --font-family: {{ .Config.FontFamily }};
      --primary-color: {{ .Config.PrimaryColor }};
      --primary-hover: {{ .Config.PrimaryHoverColor }};
      --bg-color: #f8fafc;
      --card-bg: #ffffff;
      --sidebar-bg: #ffffff;
      --text-main: #0f172a;
      --text-muted: #64748b;
      --border-color: #e2e8f0;
      --code-bg: #1e293b;
      --code-text: #f8fafc;
      --badge-bg: #f1f5f9;
      --badge-text: #334155;
      --accent-pill: #e0e7ff;
      --accent-pill-text: #3730a3;
      --sidebar-width: 280px;
    }

    [data-theme="dark"] {
      --bg-color: #0b0f17;
      --card-bg: #111827;
      --sidebar-bg: #0f172a;
      --text-main: #f3f4f6;
      --text-muted: #9ca3af;
      --border-color: #1f2937;
      --code-bg: #030712;
      --code-text: #e5e7eb;
      --badge-bg: #1f2937;
      --badge-text: #d1d5db;
      --accent-pill: #1e1b4b;
      --accent-pill-text: #a5b4fc;
    }

    * {
      box-sizing: border-box;
      margin: 0;
      padding: 0;
    }

    body {
      font-family: var(--font-family);
      background-color: var(--bg-color);
      color: var(--text-main);
      line-height: 1.5;
      -webkit-font-smoothing: antialiased;
    }

    /* Top Nav Header */
    header.top-nav {
      position: sticky;
      top: 0;
      z-index: 40;
      display: flex;
      align-items: center;
      justify-content: space-between;
      height: 64px;
      padding: 0 24px;
      background-color: var(--card-bg);
      border-bottom: 1px solid var(--border-color);
      backdrop-filter: blur(8px);
    }

    .brand {
      display: flex;
      align-items: center;
      gap: 12px;
      font-weight: 700;
      font-size: 1.125rem;
      color: var(--text-main);
      text-decoration: none;
    }

    .brand img {
      height: 28px;
    }

    .search-btn {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 8px 14px;
      background: var(--badge-bg);
      color: var(--text-muted);
      border: 1px solid var(--border-color);
      border-radius: 8px;
      font-size: 0.875rem;
      cursor: pointer;
      transition: all 0.2s ease;
    }

    .search-btn:hover {
      border-color: var(--primary-color);
      color: var(--text-main);
    }

    .search-kbd {
      padding: 2px 6px;
      background: var(--card-bg);
      border: 1px solid var(--border-color);
      border-radius: 4px;
      font-size: 0.75rem;
      font-weight: 600;
      box-shadow: 0 1px 1px rgba(0,0,0,0.05);
    }

    .header-actions {
      display: flex;
      align-items: center;
      gap: 16px;
    }

    .theme-toggle {
      background: none;
      border: none;
      cursor: pointer;
      color: var(--text-muted);
      font-size: 1.25rem;
      padding: 6px;
      border-radius: 6px;
      transition: color 0.2s;
    }

    .theme-toggle:hover {
      color: var(--text-main);
    }

    /* Container Layout */
    .app-layout {
      display: flex;
      min-height: calc(100vh - 64px);
    }

    /* Sidebar Navigation */
    aside.sidebar {
      width: var(--sidebar-width);
      flex-shrink: 0;
      background-color: var(--sidebar-bg);
      border-right: 1px solid var(--border-color);
      padding: 24px 16px;
      position: sticky;
      top: 64px;
      height: calc(100vh - 64px);
      overflow-y: auto;
    }

    .sidebar-section {
      margin-bottom: 24px;
    }

    .sidebar-title {
      font-size: 0.75rem;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      color: var(--text-muted);
      margin-bottom: 8px;
      padding: 0 8px;
    }

    .sidebar-link {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 8px 12px;
      color: var(--text-muted);
      text-decoration: none;
      font-size: 0.875rem;
      border-radius: 6px;
      transition: all 0.15s ease;
      margin-bottom: 2px;
    }

    .sidebar-link:hover, .sidebar-link.active {
      color: var(--primary-color);
      background-color: var(--badge-bg);
      font-weight: 500;
    }

    .badge-method {
      font-size: 0.65rem;
      font-weight: 700;
      padding: 2px 6px;
      border-radius: 4px;
      text-transform: uppercase;
      background: var(--accent-pill);
      color: var(--accent-pill-text);
    }

    /* Main Content Area */
    main.main-content {
      flex: 1;
      padding: 32px 48px;
      max-width: 1200px;
    }

    .doc-section {
      margin-bottom: 64px;
      scroll-margin-top: 80px;
    }

    .section-header {
      margin-bottom: 24px;
    }

    .section-title {
      font-size: 1.75rem;
      font-weight: 700;
      letter-spacing: -0.02em;
      color: var(--text-main);
      margin-bottom: 8px;
    }

    .section-meta {
      display: flex;
      align-items: center;
      gap: 12px;
      font-size: 0.875rem;
      color: var(--text-muted);
    }

    .package-tag {
      background: var(--accent-pill);
      color: var(--accent-pill-text);
      padding: 4px 10px;
      border-radius: 12px;
      font-size: 0.75rem;
      font-weight: 600;
    }

    .description {
      font-size: 1rem;
      color: var(--text-muted);
      line-height: 1.6;
      margin-top: 12px;
      white-space: pre-line;
    }

    .method-grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 32px;
      margin-top: 24px;
    }

    @media (max-width: 1024px) {
      .method-grid {
        grid-template-columns: 1fr;
      }
    }

    .card {
      background: var(--card-bg);
      border: 1px solid var(--border-color);
      border-radius: 16px;
      padding: 24px;
      box-shadow: 0 4px 20px -2px rgba(0, 0, 0, 0.03);
      transition: border-color 0.2s, box-shadow 0.2s;
    }

    .card:hover {
      border-color: var(--primary-color);
      box-shadow: 0 6px 24px -2px rgba(0, 0, 0, 0.06);
    }

    .code-box {
      background: var(--code-bg);
      color: var(--code-text);
      border-radius: 12px;
      padding: 16px;
      font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
      font-size: 0.85rem;
      overflow-x: auto;
      line-height: 1.5;
    }

    .params-table {
      width: 100%;
      border-collapse: separate;
      border-spacing: 0;
      margin-top: 16px;
    }

    .params-table th {
      text-align: left;
      font-size: 0.75rem;
      font-weight: 700;
      text-transform: uppercase;
      color: var(--text-muted);
      padding: 10px 12px;
      border-bottom: 1px solid var(--border-color);
    }

    .params-table td {
      padding: 12px;
      border-bottom: 1px solid var(--border-color);
      font-size: 0.875rem;
      vertical-align: top;
    }

    .params-table tr:last-child td {
      border-bottom: none;
    }

    .param-name {
      font-weight: 600;
      font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, monospace;
      color: var(--primary-color);
    }

    .param-type {
      font-size: 0.75rem;
      padding: 2px 8px;
      background: var(--badge-bg);
      color: var(--badge-text);
      border-radius: 6px;
      font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, monospace;
      display: inline-block;
    }

    .param-desc {
      color: var(--text-muted);
      font-size: 0.85rem;
      line-height: 1.4;
    }

    .search-modal-backdrop {
      position: fixed;
      inset: 0;
      background: rgba(15, 23, 42, 0.5);
      backdrop-filter: blur(4px);
      z-index: 100;
      display: flex;
      align-items: flex-start;
      justify-content: center;
      padding-top: 10vh;
      opacity: 0;
      pointer-events: none;
      transition: opacity 0.2s ease;
    }

    .search-modal-backdrop.open {
      opacity: 1;
      pointer-events: auto;
    }

    .search-modal {
      background: var(--card-bg);
      border: 1px solid var(--border-color);
      width: 100%;
      max-width: 640px;
      border-radius: 16px;
      box-shadow: 0 20px 25px -5px rgba(0,0,0,0.2);
      overflow: hidden;
    }

    .search-input-wrapper {
      display: flex;
      align-items: center;
      padding: 16px 20px;
      border-bottom: 1px solid var(--border-color);
    }

    .search-input {
      width: 100%;
      background: none;
      border: none;
      outline: none;
      font-size: 1rem;
      color: var(--text-main);
      font-family: inherit;
    }

    .search-results {
      max-height: 400px;
      overflow-y: auto;
      padding: 12px;
    }

    .search-result-item {
      display: flex;
      flex-direction: column;
      padding: 12px;
      border-radius: 8px;
      cursor: pointer;
      text-decoration: none;
      color: inherit;
      transition: background 0.15s;
    }

    .search-result-item:hover, .search-result-item.selected {
      background: var(--badge-bg);
    }

    .search-result-title {
      font-weight: 600;
      font-size: 0.9rem;
      display: flex;
      align-items: center;
      justify-content: space-between;
    }

    .search-result-desc {
      font-size: 0.8rem;
      color: var(--text-muted);
      margin-top: 4px;
    }
  </style>
</head>
<body>
  <header class="top-nav">
    <a href="#" class="brand">
      {{ if .Config.LogoURL }}
        <img src="{{ .Config.LogoURL }}" alt="Logo" />
      {{ end }}
      <span>{{ .Config.Title }}</span>
    </a>

    <button class="search-btn" id="search-trigger" onclick="toggleSearchModal(true)">
      <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
      <span>Search documentation...</span>
      <span class="search-kbd">⌘K</span>
    </button>

    <div class="header-actions">
      <button class="theme-toggle" onclick="toggleTheme()" title="Toggle dark/light mode">
        <span id="theme-icon">🌙</span>
      </button>
    </div>
  </header>

  <div class="app-layout">
    <aside class="sidebar">
      {{ range .DocData.Files }}
        {{ if .Services }}
          <div class="sidebar-section">
            <div class="sidebar-title">Services</div>
            {{ range $svc := .Services }}
              <a href="#service-{{ $svc.FullName }}" class="sidebar-link">
                <span>{{ $svc.Name }}</span>
              </a>
              {{ range $svc.Methods }}
                <a href="#method-{{ $svc.FullName }}-{{ .Name }}" class="sidebar-link" style="padding-left: 24px;">
                  <span>{{ .Name }}</span>
                  <span class="badge-method">RPC</span>
                </a>
              {{ end }}
            {{ end }}
          </div>
        {{ end }}

        {{ if .Messages }}
          <div class="sidebar-section">
            <div class="sidebar-title">Messages</div>
            {{ range .Messages }}
              {{ if not .IsMapEntry }}
                <a href="#message-{{ .FullName }}" class="sidebar-link">
                  <span>{{ .Name }}</span>
                </a>
              {{ end }}
            {{ end }}
          </div>
        {{ end }}

        {{ if .Enums }}
          <div class="sidebar-section">
            <div class="sidebar-title">Enums</div>
            {{ range .Enums }}
              <a href="#enum-{{ .FullName }}" class="sidebar-link">
                <span>{{ .Name }}</span>
              </a>
            {{ end }}
          </div>
        {{ end }}

        {{ if .Extensions }}
          <div class="sidebar-section">
            <div class="sidebar-title">Extensions</div>
            {{ range .Extensions }}
              <a href="#extension-{{ .Name }}" class="sidebar-link">
                <span>{{ .Name }}</span>
              </a>
            {{ end }}
          </div>
        {{ end }}
      {{ end }}
    </aside>

    <main class="main-content">
      {{ range $file := .DocData.Files }}
        {{ range $service := $file.Services }}
          <section id="service-{{ $service.FullName }}" class="doc-section">
            <div class="section-header">
              <div class="section-meta">
                <span class="package-tag">{{ $file.Package }}</span>
                <span>Service</span>
              </div>
              <h1 class="section-title">{{ $service.Name }}</h1>
              {{ if $service.Comments.Description }}
                <div class="description">{{ $service.Comments.Description }}</div>
              {{ end }}
            </div>

            {{ range $method := $service.Methods }}
              <div id="method-{{ $service.FullName }}-{{ $method.Name }}" class="doc-section" style="margin-top: 40px;">
                <h2>{{ $method.Name }}</h2>
                {{ if $method.Comments.Description }}
                  <p class="description">{{ $method.Comments.Description }}</p>
                {{ end }}

                <div class="method-grid">
                  <div class="card">
                    <h3>Request & Response</h3>
                    <table class="params-table">
                      <thead>
                        <tr>
                          <th>Type</th>
                          <th>Schema</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr>
                          <td>Input Request</td>
                          <td><a href="#message-{{ $method.InputType }}" class="param-name">{{ $method.InputType }}</a></td>
                        </tr>
                        <tr>
                          <td>Output Response</td>
                          <td><a href="#message-{{ $method.OutputType }}" class="param-name">{{ $method.OutputType }}</a></td>
                        </tr>
                      </tbody>
                    </table>
                  </div>

                  <div class="card">
                    <h3>gRPC Definition</h3>
                    <div class="code-box">rpc {{ $method.Name }} ({{ $method.InputType }}) returns ({{ $method.OutputType }});</div>
                  </div>
                </div>
              </div>
            {{ end }}
          </section>
        {{ end }}

        {{ if $file.Messages }}
          <section class="doc-section">
            <h1 class="section-title">Messages</h1>
            {{ range $msg := $file.Messages }}
              {{ if not $msg.IsMapEntry }}
                <div id="message-{{ $msg.FullName }}" class="card" style="margin-bottom: 24px;">
                  <div class="section-meta">
                    <span class="package-tag">{{ $file.Package }}</span>
                    <span>Message</span>
                  </div>
                  <h2 style="margin-top: 8px;">{{ $msg.Name }}</h2>
                  {{ if $msg.Comments.Description }}
                    <p class="description">{{ $msg.Comments.Description }}</p>
                  {{ end }}

                  {{ if $msg.Fields }}
                    <table class="params-table">
                      <thead>
                        <tr>
                          <th>Field</th>
                          <th>Type</th>
                          <th>Rule</th>
                          <th>Description</th>
                        </tr>
                      </thead>
                      <tbody>
                        {{ range $field := $msg.Fields }}
                          <tr>
                            <td class="param-name">{{ $field.Name }}</td>
                            <td>
                              {{ if $field.TypeName }}
                                {{ if eq $field.Type "enum" }}
                                  <a href="#enum-{{ $field.TypeName }}" class="param-type">{{ $field.TypeDisplay }}</a>
                                {{ else }}
                                  <a href="#message-{{ $field.TypeName }}" class="param-type">{{ $field.TypeDisplay }}</a>
                                {{ end }}
                              {{ else }}
                                <span class="param-type">{{ $field.Type }}</span>
                              {{ end }}
                            </td>
                            <td><span class="param-desc">{{ $field.Label }}</span></td>
                            <td class="param-desc">
                              {{ $field.Comments.Description }}
                              {{ if $field.OneofName }}
                                <div style="margin-top: 4px;"><span class="badge-method">oneof {{ $field.OneofName }}</span></div>
                              {{ end }}
                            </td>
                          </tr>
                        {{ end }}
                      </tbody>
                    </table>
                  {{ end }}
                </div>
              {{ end }}
            {{ end }}
          </section>
        {{ end }}

        {{ if $file.Enums }}
          <section class="doc-section">
            <h1 class="section-title">Enums</h1>
            {{ range $enum := $file.Enums }}
              <div id="enum-{{ $enum.FullName }}" class="card" style="margin-bottom: 24px;">
                <div class="section-meta">
                  <span class="package-tag">{{ $file.Package }}</span>
                  <span>Enum</span>
                </div>
                <h2 style="margin-top: 8px;">{{ $enum.Name }}</h2>
                {{ if $enum.Comments.Description }}
                  <p class="description">{{ $enum.Comments.Description }}</p>
                {{ end }}

                <table class="params-table">
                  <thead>
                    <tr>
                      <th>Name</th>
                      <th>Number</th>
                      <th>Description</th>
                    </tr>
                  </thead>
                  <tbody>
                    {{ range $val := $enum.Values }}
                      <tr>
                        <td class="param-name">{{ $val.Name }}</td>
                        <td><span class="param-type">{{ $val.Number }}</span></td>
                        <td class="param-desc">{{ $val.Comments.Description }}</td>
                      </tr>
                    {{ end }}
                  </tbody>
                </table>
              </div>
            {{ end }}
          </section>
        {{ end }}

        {{ if $file.Extensions }}
          <section class="doc-section">
            <h1 class="section-title">Extensions</h1>
            {{ range $ext := $file.Extensions }}
              <div id="extension-{{ $ext.Name }}" class="card" style="margin-bottom: 24px;">
                <div class="section-meta">
                  <span class="package-tag">{{ $file.Package }}</span>
                  <span>Extension</span>
                </div>
                <h2 style="margin-top: 8px;">{{ $ext.Name }}</h2>
                <p class="param-desc">Extends <code>{{ $ext.Extendee }}</code> with type <code>{{ $ext.Type }}</code> (Tag: {{ $ext.Number }})</p>
                {{ if $ext.Comments.Description }}
                  <p class="description">{{ $ext.Comments.Description }}</p>
                {{ end }}
              </div>
            {{ end }}
          </section>
        {{ end }}
      {{ end }}
    </main>
  </div>

  <div class="search-modal-backdrop" id="search-modal" onclick="if(event.target===this) toggleSearchModal(false)">
    <div class="search-modal">
      <div class="search-input-wrapper">
        <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24" style="margin-right: 12px; color: var(--text-muted);"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
        <input type="text" id="search-input" class="search-input" placeholder="Search services, methods, messages, fields..." oninput="handleSearch(this.value)" />
      </div>
      <div class="search-results" id="search-results">
        <div style="padding: 16px; text-align: center; color: var(--text-muted); font-size: 0.875rem;">Type to search...</div>
      </div>
    </div>
  </div>

  <script>
    const docData = {{ .DocDataJSON }};

    const searchIndex = [];
    if (docData && docData.files) {
      docData.files.forEach(f => {
        (f.services || []).forEach(s => {
          searchIndex.push({ type: 'Service', name: s.name, desc: s.comments.description, link: '#service-' + s.full_name });
          (s.methods || []).forEach(m => {
            searchIndex.push({ type: 'RPC Method', name: s.name + '.' + m.name, desc: m.comments.description, link: '#method-' + s.full_name + '-' + m.name });
          });
        });
        (f.messages || []).forEach(m => {
          if (!m.is_map_entry) {
            searchIndex.push({ type: 'Message', name: m.name, desc: m.comments.description, link: '#message-' + m.full_name });
            (m.fields || []).forEach(fld => {
              searchIndex.push({ type: 'Field', name: m.name + '.' + fld.name, desc: fld.comments.description, link: '#message-' + m.full_name });
            });
          }
        });
        (f.enums || []).forEach(e => {
          searchIndex.push({ type: 'Enum', name: e.name, desc: e.comments.description, link: '#enum-' + e.full_name });
        });
      });
    }

    function toggleSearchModal(show) {
      const modal = document.getElementById('search-modal');
      const input = document.getElementById('search-input');
      if (show) {
        modal.classList.add('open');
        input.value = '';
        input.focus();
        handleSearch('');
      } else {
        modal.classList.remove('open');
      }
    }

    function handleSearch(query) {
      const container = document.getElementById('search-results');
      query = query.trim().toLowerCase();
      if (!query) {
        container.innerHTML = '<div style="padding: 16px; text-align: center; color: var(--text-muted); font-size: 0.875rem;">Type to search...</div>';
        return;
      }

      const results = searchIndex.filter(item =>
        item.name.toLowerCase().includes(query) || (item.desc && item.desc.toLowerCase().includes(query))
      );

      if (results.length === 0) {
        container.innerHTML = '<div style="padding: 16px; text-align: center; color: var(--text-muted); font-size: 0.875rem;">No results found</div>';
        return;
      }

      container.innerHTML = results.slice(0, 10).map(item =>
        '<a href="' + item.link + '" class="search-result-item" onclick="toggleSearchModal(false)">' +
          '<div class="search-result-title">' +
            '<span>' + item.name + '</span>' +
            '<span class="badge-method">' + item.type + '</span>' +
          '</div>' +
          (item.desc ? '<div class="search-result-desc">' + item.desc.slice(0, 100) + '...</div>' : '') +
        '</a>'
      ).join('');
    }

    window.addEventListener('keydown', (e) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault();
        const modal = document.getElementById('search-modal');
        toggleSearchModal(!modal.classList.contains('open'));
      }
      if (e.key === 'Escape') {
        toggleSearchModal(false);
      }
    });

    function toggleTheme() {
      const current = document.documentElement.getAttribute('data-theme');
      const next = current === 'dark' ? 'light' : 'dark';
      document.documentElement.setAttribute('data-theme', next);
      document.getElementById('theme-icon').textContent = next === 'dark' ? '☀️' : '🌙';
    }
  </script>
</body>
</html>`
