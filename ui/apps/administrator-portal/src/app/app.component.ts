import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { SnackbarComponent } from '@shikshakul/shared/ui/snackbar';

@Component({
  imports: [RouterModule, SnackbarComponent],
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
})
export class App {
  protected title = 'administrator-portal';
}
