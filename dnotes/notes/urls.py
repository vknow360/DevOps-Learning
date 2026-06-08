from django.urls import path
from . import views

urlpatterns = [
    path("", view=views.note_list, name="notes"),
    path("new/", view=views.note_create, name="note_create"),
    path("<int:note_id>/", view=views.note_detail, name="note"),
    path("<int:note_id>/edit/", view=views.note_edit, name="note_edit"),
    path("<int:note_id>/delete/", view=views.note_delete, name="note_delete"),
]