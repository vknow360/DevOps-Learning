from django.forms.models import ModelForm
from django.forms import Textarea
from .models import Note


class NoteForm(ModelForm):
    class Meta:
        model = Note
        fields = ["message"]
        widgets = {
            'message': Textarea(attrs={
                'class': 'form-control-textarea',
                'placeholder': 'Start typing your note here...',
                'rows': 12,
            })
        }

    def save(self, commit=True):
        note = super().save(commit=False)
        note.title = self.cleaned_data["message"][:10]
        if commit:
            note.save()
        return note