import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { SnackbarService } from '../service/snackbar.service';

@Component({
  selector: 'shikshakul-snackbar',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './snackbar.component.html',
  styleUrl: './snackbar.component.scss',
})
export class SnackbarComponent {
  public snackbarService = inject(SnackbarService);
  activeSnackbar$ = this.snackbarService.snackbar$;

  close() {
    this.snackbarService.hide();
  }
}
