import { CommonModule } from '@angular/common';
import { Component, signal } from '@angular/core';
import { FeeSetupComponent } from '../components/fee-setup/fee-setup.component';
import { FeeCollectionComponent } from '../components/fee-collection/fee-collection.component';

@Component({
  selector: 'shikshakul-fee-management-page',
  standalone: true,
  imports: [CommonModule, FeeSetupComponent, FeeCollectionComponent],
  templateUrl: './fee-management-page.component.html',
  styleUrl: './fee-management-page.component.scss',
})
export class FeeManagementPageComponent {
  activeTab: 'COLLECTION' | 'SETUP' = 'COLLECTION';
  loading = signal(false);

  setLoading(state: boolean) {
    this.loading.set(state);
  }
}
