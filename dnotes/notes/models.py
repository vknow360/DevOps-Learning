from django.db import models
from django.utils import timezone

# Create your models here.
class Note(models.Model):
    title = models.CharField(max_length=50)
    message = models.TextField()
    created = models.DateTimeField(default=timezone.now)

    def __str__(self):
        return self.title