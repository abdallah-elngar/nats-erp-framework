{{define "content"}}
<div class="p-6">
    <div class="flex justify-between items-center mb-4">
        <h2 class="text-xl font-bold">{{.Title}} Details</h2>
        <button class="text-gray-500 hover:text-gray-700" @click="showModal = false">
            <i class="fas fa-times"></i>
        </button>
    </div>
    
    <div class="bg-gray-50 rounded-lg p-4">
        <dl class="grid grid-cols-2 gap-4">
            <div class="col-span-1">
                <dt class="text-sm font-medium text-gray-500">ID</dt>
                <dd class="text-sm text-gray-900">{{.Item.ID}}</dd>
            </div>
            {{range .Fields}}
            <div class="col-span-1">
                <dt class="text-sm font-medium text-gray-500">{{.DisplayName}}</dt>
                <dd class="text-sm text-gray-900">
                    {{if eq .Type "boolean"}}
                    <span class="badge {{if index $.Item .Name}}badge-success{{else}}badge-danger{{end}}">
                        {{if index $.Item .Name}}Yes{{else}}No{{end}}
                    </span>
                    {{else if eq .Type "datetime"}}
                    {{formatDate (index $.Item .Name)}}
                    {{else if eq .Type "date"}}
                    {{formatDate (index $.Item .Name)}}
                    {{else if eq .Type "number"}}
                    {{formatNumber (index $.Item .Name)}}
                    {{else}}
                    {{index $.Item .Name}}
                    {{end}}
                </dd>
            </div>
            {{end}}
            <div class="col-span-1">
                <dt class="text-sm font-medium text-gray-500">Created At</dt>
                <dd class="text-sm text-gray-900">{{formatDate .Item.CreatedAt}}</dd>
            </div>
            <div class="col-span-1">
                <dt class="text-sm font-medium text-gray-500">Updated At</dt>
                <dd class="text-sm text-gray-900">{{formatDate .Item.UpdatedAt}}</dd>
            </div>
        </dl>
    </div>
    
    <div class="flex justify-end gap-2 mt-4">
        <button class="btn btn-secondary" @click="showModal = false">
            Close
        </button>
        <button class="btn btn-primary"
                hx-get="/{{.AppName}}/{{.ModelLower}}/{{.Item.ID}}/edit"
                hx-target="#modal-content"
                hx-trigger="click"
                @click="showModal = true">
            <i class="fas fa-edit mr-2"></i>Edit
        </button>
    </div>
</div>
{{end}}