{{define "content"}}
<div class="p-6">
    <div class="flex justify-between items-center mb-4">
        <h2 class="text-xl font-bold">Create New {{.Title}}</h2>
        <button class="text-gray-500 hover:text-gray-700" @click="showModal = false">
            <i class="fas fa-times"></i>
        </button>
    </div>
    
    <form hx-post="/{{.AppName}}/{{.ModelLower}}"
          hx-target="#modal-content"
          hx-swap="innerHTML"
          hx-trigger="submit"
          data-validate>
        
        {{range .Fields}}
        <div class="form-group">
            <label class="form-label" for="{{.Name}}">{{.DisplayName}}</label>
            {{if eq .Type "textarea"}}
            <textarea class="form-control" id="{{.Name}}" name="{{.Name}}" rows="3" {{if .Required}}required{{end}}></textarea>
            {{else if eq .Type "select"}}
            <select class="form-control" id="{{.Name}}" name="{{.Name}}" {{if .Required}}required{{end}}>
                <option value="">Select {{.DisplayName}}</option>
                {{range .Options}}
                <option value="{{.Value}}">{{.Label}}</option>
                {{end}}
            </select>
            {{else if eq .Type "checkbox"}}
            <input type="checkbox" class="form-checkbox" id="{{.Name}}" name="{{.Name}}" value="1">
            {{else if eq .Type "date"}}
            <input type="date" class="form-control" id="{{.Name}}" name="{{.Name}}" {{if .Required}}required{{end}}>
            {{else if eq .Type "datetime"}}
            <input type="datetime-local" class="form-control" id="{{.Name}}" name="{{.Name}}" {{if .Required}}required{{end}}>
            {{else if eq .Type "email"}}
            <input type="email" class="form-control" id="{{.Name}}" name="{{.Name}}" {{if .Required}}required{{end}}>
            {{else if eq .Type "password"}}
            <input type="password" class="form-control" id="{{.Name}}" name="{{.Name}}" {{if .Required}}required{{end}}>
            {{else if eq .Type "number"}}
            <input type="number" class="form-control" id="{{.Name}}" name="{{.Name}}" {{if .Required}}required{{end}}>
            {{else}}
            <input type="text" class="form-control" id="{{.Name}}" name="{{.Name}}" {{if .Required}}required{{end}}>
            {{end}}
            {{if .Help}}
            <small class="text-gray-500 text-sm">{{.Help}}</small>
            {{end}}
        </div>
        {{end}}
        
        <div class="flex justify-end gap-2 mt-4">
            <button type="button" class="btn btn-secondary" @click="showModal = false">
                Cancel
            </button>
            <button type="submit" class="btn btn-primary">
                <i class="fas fa-save mr-2"></i>Save
            </button>
        </div>
    </form>
</div>
{{end}}