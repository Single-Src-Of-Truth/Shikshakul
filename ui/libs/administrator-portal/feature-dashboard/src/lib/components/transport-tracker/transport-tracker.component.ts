import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';

@Component({
  selector: 'shikshakul-transport-tracker',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './transport-tracker.component.html',
  styleUrl: './transport-tracker.component.scss',
})
export class TransportTrackerComponent {
  routes = [
    {
      name: 'Route 12 - South City',
      status: 'On Time',
      detail: 'Arrived at Stop 4',
      vehicleNo: 'WB-04-1234',
      state: 'success', // Green
    },
    {
      name: 'Route 08 - Lake Gardens',
      status: 'Delayed',
      detail: 'by 10 mins • Traffic',
      vehicleNo: 'WB-04-5678',
      state: 'warning', // Orange
    },
    {
      name: 'Route 03 - Salt Lake',
      status: 'Trip Completed',
      detail: '',
      vehicleNo: 'WB-04-9012',
      state: 'neutral', // Grey
    },
  ];
}
