from django.shortcuts import render, redirect, get_object_or_404
from .forms import NoteForm
from .models import Note

# Create your views here.
def note_list(request):
    notes = Note.objects.all().order_by("-created")
    return render(request, "notes/list.html", {'notes': notes})

def note_detail(request, note_id):
    note = get_object_or_404(Note, id=note_id)
    return render(request, "notes/detail.html", {'note': note})

def note_create(request):
    form = NoteForm(request.POST or None)
    if form.is_valid():
        form.save()
        return redirect("notes")
    return render(request, "notes/create.html", {'form' : form})

def note_edit(request, note_id):
    note = get_object_or_404(Note, id=note_id)
    form = NoteForm(request.POST or None, instance=note)
    if form.is_valid():
        form.save()
        return redirect("note", note_id=note.id)
    return render(request, "notes/edit.html", {'form': form, 'note': note})

def note_delete(request, note_id):
    note = get_object_or_404(Note, id=note_id)
    if request.method == "POST":
        note.delete()
        return redirect("notes")
    return render(request, "notes/delete_confirm.html", {'note': note})