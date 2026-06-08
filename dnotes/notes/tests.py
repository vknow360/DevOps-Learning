from django.test import TestCase
from django.urls import reverse
from .models import Note


class NoteModelTest(TestCase):
    def test_note_creation_and_title_auto_extraction(self):
        # Test model creation and form saving behavior
        # Note: the form save auto-titles, but let's test model creation directly
        note = Note.objects.create(title="Short title", message="This is a test message content.")
        self.assertEqual(note.title, "Short title")
        self.assertEqual(note.message, "This is a test message content.")
        self.assertIsNotNone(note.created)
        self.assertEqual(str(note), "Short title")


class NoteViewsTest(TestCase):
    def setUp(self):
        self.note = Note.objects.create(title="Test Note", message="This is a test note.")

    def test_note_list_view(self):
        response = self.client.get(reverse("notes"))
        self.assertEqual(response.status_code, 200)
        self.assertTemplateUsed(response, "notes/list.html")
        self.assertContains(response, "Test Note")

    def test_note_detail_view(self):
        response = self.client.get(reverse("note", args=[self.note.id]))
        self.assertEqual(response.status_code, 200)
        self.assertTemplateUsed(response, "notes/detail.html")
        self.assertContains(response, "This is a test note.")

    def test_note_create_view_get(self):
        response = self.client.get(reverse("note_create"))
        self.assertEqual(response.status_code, 200)
        self.assertTemplateUsed(response, "notes/create.html")

    def test_note_create_view_post(self):
        message_text = "New note content written here."
        response = self.client.post(reverse("note_create"), data={"message": message_text})
        # Redirect to notes list view
        self.assertRedirects(response, reverse("notes"))
        
        # Verify note was created and title auto-extracted
        new_note = Note.objects.get(message=message_text)
        self.assertEqual(new_note.title, "New note c") # first 10 characters

    def test_note_edit_view_get(self):
        response = self.client.get(reverse("note_edit", args=[self.note.id]))
        self.assertEqual(response.status_code, 200)
        self.assertTemplateUsed(response, "notes/edit.html")

    def test_note_edit_view_post(self):
        updated_message = "Updated message content."
        response = self.client.post(reverse("note_edit", args=[self.note.id]), data={"message": updated_message})
        # Redirects to detail view
        self.assertRedirects(response, reverse("note", args=[self.note.id]))
        
        # Verify model update
        self.note.refresh_from_db()
        self.assertEqual(self.note.message, updated_message)
        self.assertEqual(self.note.title, "Updated me") # first 10 characters

    def test_note_delete_view_get(self):
        response = self.client.get(reverse("note_delete", args=[self.note.id]))
        self.assertEqual(response.status_code, 200)
        self.assertTemplateUsed(response, "notes/delete_confirm.html")

    def test_note_delete_view_post(self):
        response = self.client.post(reverse("note_delete", args=[self.note.id]))
        self.assertRedirects(response, reverse("notes"))
        # Verify deletion
        self.assertFalse(Note.objects.filter(id=self.note.id).exists())
