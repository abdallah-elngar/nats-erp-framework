{{define "content"}}
<div>
    <div class="flex justify-between items-center mb-6">
        <h1 class="text-2xl font-bold text-gray-800">{{.Title}}</h1>
        <button class="btn btn-primary"
                hx-get="/{{.AppName}}/{{.ModelLower}}/create"
                hx-target="#modal-content"
                hx-trigger="click"
                @click="showModal = true">
            <i class="fas fa-plus mr-2"></i>Add New
        </button>
    </div>
    
    <!-- Search -->
    <div class="mb-4">
        <div class="relative">
            <span class="absolute inset-y-0 left-0 pl-3 flex items-center text-gray-500">
                <i class="fas fa-search"></i>
            </span>
            <input type="text" 
                   placeholder="Search..."
                   class="w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
                   hx-get="/{{.AppName}}/{{.ModelLower}}/search"
                   hx-trigger="keyup changed delay:300ms"
                   hx-target="#table-content"
                   hx-swap="outerHTML">
        </div>
    </div>
    
    <!-- Table -->
    <div class="bg-white rounded-lg shadow overflow-hidden">
        <div id="table-content">
            <table class="w-full">
                <thead class="bg-gray-50">
                    <tr>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
                        {{range .Fields}}
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{.Name}}</th>
                        {{end}}
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-gray-200">
                    {{range .Items}}
                    <tr>
                        <td class="px-6 py-4">{{.ID}}</td>
                        {{range $.Fields}}
                        <td class="px-6 py-4">{{index $ .Name}}</td>
                        {{end}}
                        <td class="px-6 py-4">
                            <button class="text-blue-600 hover:text-blue-800 mr-2"
                                    hx-get="/{{$.AppName}}/{{$.ModelLower}}/{{.ID}}"
                                    hx-target="#modal-content"
                                    @click="showModal = true">
                                <i class="fas fa-eye"></i>
                            </button>
                            <button class="text-green-600 hover:text-green-800 mr-2"
                                    hx-get="/{{$.AppName}}/{{$.ModelLower}}/{{.ID}}/edit"
                                    hx-target="#modal-content"
                                    @click="showModal = true">
                                <i class="fas fa-edit"></i>
                            </button>
                            <button class="text-red-600 hover:text-red-800"
                                    hx-delete="/{{$.AppName}}/{{$.ModelLower}}/{{.ID}}"
                                    hx-confirm="Are you sure?"
                                    hx-target="#table-content"
                                    hx-swap="outerHTML">
                                <i class="fas fa-trash"></i>
                            </button>
                        </td>
                    </tr>
                    {{else}}
                    <tr>
                        <td colspan="{{len .Fields | add 2}}" class="px-6 py-4 text-center text-gray-500">
                            No records found
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
        
        <!-- Pagination -->
        {{if .Pagination}}
        <div class="px-6 py-3 bg-gray-50 border-t flex items-center justify-between">
            <div class="text-sm text-gray-700">
                Showing {{.Pagination.Start}} to {{.Pagination.End}} of {{.Pagination.Total}}
            </div>
            <div class="flex gap-2">
                {{if .Pagination.PrevPage}}
                <button hx-get="/{{.AppName}}/{{.ModelLower}}?page={{.Pagination.PrevPage}}"
                        hx-target="#table-content"
                        hx-swap="outerHTML"
                        class="btn btn-secondary btn-sm">
                    Previous
                </button>
                {{end}}
                {{range .Pagination.Pages}}
                <button hx-get="/{{$.AppName}}/{{$.ModelLower}}?page={{.}}"
                        hx-target="#table-content"
                        hx-swap="outerHTML"
                        class="btn {{if eq . $.Pagination.Current}}btn-primary{{else}}btn-secondary{{end}} btn-sm">
                    {{.}}
                </button>
                {{end}}
                {{if .Pagination.NextPage}}
                <button hx-get="/{{.AppName}}/{{.ModelLower}}?page={{.Pagination.NextPage}}"
                        hx-target="#table-content"
                        hx-swap="outerHTML"
                        class="btn btn-secondary btn-sm">
                    Next
                </button>
                {{end}}
            </div>
        </div>
        {{end}}
    </div>
</div>
{{end}}