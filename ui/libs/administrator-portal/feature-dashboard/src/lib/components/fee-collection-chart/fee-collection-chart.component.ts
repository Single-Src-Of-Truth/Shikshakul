import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { ChartConfiguration, ChartData, ChartType } from 'chart.js';
import { BaseChartDirective } from 'ng2-charts';

@Component({
  selector: 'shikshakul-fee-collection-chart',
  standalone: true,
  imports: [CommonModule, BaseChartDirective],
  templateUrl: './fee-collection-chart.component.html',
  styleUrl: './fee-collection-chart.component.scss',
})
export class FeeCollectionChartComponent {
  public lineChartType: ChartType = 'line';

  public lineChartData: ChartData<'line'> = {
    labels: ['Week 1', 'Week 2', 'Week 3', 'Week 4', 'Week 5'],
    datasets: [
      {
        data: [20, 35, 25, 65, 60], // The curve data points
        label: 'Collections',
        backgroundColor: (context) => {
          // Creating a Gradient Fill
          const ctx = context.chart.ctx;
          const gradient = ctx.createLinearGradient(0, 0, 0, 400);
          gradient.addColorStop(0, 'rgba(24, 144, 255, 0.2)'); // Light Blue at top
          gradient.addColorStop(1, 'rgba(24, 144, 255, 0.0)'); // Transparent at bottom
          return gradient;
        },
        borderColor: '#1890FF', // The solid blue line
        pointBackgroundColor: '#fff',
        pointBorderColor: '#1890FF',
        pointHoverBackgroundColor: '#1890FF',
        pointHoverBorderColor: '#fff',
        fill: true, // Fill the area under the line
        tension: 0.4, // This creates the "Smooth Curve" (Bezier)
      },
    ],
  };

  public lineChartOptions: ChartConfiguration['options'] = {
    responsive: true,
    maintainAspectRatio: false,
    elements: {
      point: {
        radius: 0, // Hide points by default for a clean look
        hitRadius: 10,
        hoverRadius: 6,
      },
    },
    plugins: {
      legend: { display: false }, // Hide legend (Title handles it)
      tooltip: {
        backgroundColor: '#fff',
        titleColor: '#2B3674',
        bodyColor: '#1890FF',
        borderColor: '#E0E5F2',
        borderWidth: 1,
        displayColors: false,
      },
    },
    scales: {
      x: {
        grid: { display: false },
        ticks: { color: '#A3AED0', font: { family: "'DM Sans', sans-serif" } },
        border: { display: false },
      },
      y: {
        display: false, // Completely hide Y-axis for that clean "minimal" look
      },
    },
  };
}
