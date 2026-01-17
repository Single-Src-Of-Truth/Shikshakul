import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';

@Component({
  selector: 'shikshakul-dashboard-stats',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './dashboard-stats.component.html',
  styleUrl: './dashboard-stats.component.scss',
})
export class DashboardStatsComponent {
  stats = [
    {
      label: 'Total Students',
      value: '2,450',
      change: '+1.2%',
      trend: 'up', // 'up' | 'down' | 'neutral'
      trendLabel: 'vs last month',
      icon: 'school',
      iconColor: '#1890FF',
      iconBg: '#E6F7FF',
    },
    {
      label: 'Total Teachers',
      value: '112',
      change: 'No change',
      trend: 'neutral',
      trendLabel: '',
      icon: 'badge',
      iconColor: '#FAAD14',
      iconBg: '#FFF7E6',
    },
    {
      label: 'Active Classes',
      value: '48',
      change: '',
      trend: 'neutral',
      trendLabel: 'Academic Year 23-24',
      icon: 'class',
      iconColor: '#1890FF',
      iconBg: '#E6F7FF',
    },
    {
      label: "Today's Attendance",
      value: '94%',
      change: '-2%',
      trend: 'down',
      trendLabel: 'vs yesterday',
      icon: 'check_circle',
      iconColor: '#52C41A',
      iconBg: '#F6FFED',
    },
  ];
}
