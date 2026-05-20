import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import {
  Chart,
  ChartConfiguration,
  ChartData,
  ChartType,
  registerables,
} from 'chart.js';
import { BaseChartDirective } from 'ng2-charts';

Chart.register(...registerables);

@Component({
  selector: 'shikshakul-attendance-chart',
  standalone: true,
  imports: [CommonModule, BaseChartDirective],
  templateUrl: './attendance-chart.component.html',
  styleUrl: './attendance-chart.component.scss',
})
export class AttendanceChartComponent {
  public barChartType: ChartType = 'bar';

  public barChartData: ChartData<'bar'> = {
    labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'],
    datasets: [
      {
        data: [95, 92, 98, 94, 90, 85],
        label: 'Students',
        backgroundColor: '#1A73E8', // <--- NEW ROYAL BLUE
        hoverBackgroundColor: '#174EA6', // Darker on hover
        barThickness: 16, // Slightly thicker bars
        borderRadius: 2, // Less rounded, more "data" feel
      },
      {
        data: [100, 98, 99, 98, 97, 95],
        label: 'Staff',
        backgroundColor: '#E8F0FE', // <--- Light Tint of Royal Blue
        hoverBackgroundColor: '#D2E3FC',
        barThickness: 16,
        borderRadius: 2,
      },
    ],
  };

  public barChartOptions: ChartConfiguration['options'] = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'bottom',
        align: 'end', // Moves legend to the right like in Figma
        labels: {
          usePointStyle: true, // Makes the legend icon a circle
          pointStyle: 'circle',
          padding: 20,
          color: '#A3AED0', // Grey text
          font: { family: "'DM Sans', sans-serif", size: 12 },
        },
      },
      tooltip: {
        backgroundColor: '#fff',
        titleColor: '#2B3674',
        bodyColor: '#A3AED0',
        borderColor: '#E0E5F2',
        borderWidth: 1,
        padding: 10,
        displayColors: true,
        usePointStyle: true,
      },
    },
    scales: {
      x: {
        grid: { display: false }, // No vertical grid lines
        ticks: { color: '#A3AED0', font: { family: "'DM Sans', sans-serif" } },
        border: { display: false }, // Remove axis line
      },
      y: {
        grid: { display: false }, // Clean look (no horizontal grid lines)
        ticks: { display: false }, // Hide Y-axis numbers
        border: { display: false },
      },
    },
  };
}
