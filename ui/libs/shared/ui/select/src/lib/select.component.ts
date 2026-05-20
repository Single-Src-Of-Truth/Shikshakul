import {
  Component,
  ElementRef,
  EventEmitter,
  forwardRef,
  HostListener,
  Input,
  Output,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';

export interface SelectOption {
  label: string;
  value: string;
}

@Component({
  selector: 'shikshakul-select',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './select.component.html',
  styleUrl: './select.component.scss',
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => SelectComponent),
      multi: true,
    },
  ],
})
export class SelectComponent implements ControlValueAccessor {
  @Input() options: SelectOption[] = [];
  @Input() placeholder = 'Select an option';
  @Input() disabled = false;
  @Input() hasError = false;
  @Output() onSelect = new EventEmitter<string>();

  isOpen = false;
  openUpwards = false;
  value: string | null = null;
  selectedLabel: string | null = null;

  // ControlValueAccessor functions
  onChange: any = () => {
    /* Empty function for ControlValueAccessor */
  };
  onTouch: any = () => {
    /* Empty function for ControlValueAccessor */
  };

  constructor(private elementRef: ElementRef) { }

  @HostListener('document:click', ['$event'])
  onClickOutside(event: Event) {
    if (!this.elementRef.nativeElement.contains(event.target)) {
      this.closeDropdown();
    }
  }

  @HostListener('window:resize')
  @HostListener('window:scroll')
  onWindowChange() {
    if (this.isOpen) {
      this.checkPosition();
    }
  }

  toggleDropdown() {
    if (this.disabled) return;
    this.isOpen = !this.isOpen;
    if (this.isOpen) {
      this.checkPosition();
      this.scrollToSelected();
    } else {
      this.onTouch();
    }
  }

  private checkPosition() {
    const rect = this.elementRef.nativeElement.getBoundingClientRect();
    const viewportHeight = window.innerHeight;
    const spaceBelow = viewportHeight - rect.bottom;
    const spaceAbove = rect.top;
    const dropdownHeight = 250;

    this.openUpwards = spaceBelow < dropdownHeight && spaceAbove > spaceBelow;
  }

  private scrollToSelected() {
    setTimeout(() => {
      const dropdown = this.elementRef.nativeElement.querySelector('.skl-select-dropdown');
      const selectedOption = dropdown?.querySelector('.skl-select-option.selected') as HTMLElement;

      if (dropdown && selectedOption) {
        const scrollAmount = selectedOption.offsetTop - (dropdown.clientHeight / 2) + (selectedOption.clientHeight / 2);
        dropdown.scrollTo({ top: scrollAmount, behavior: 'instant' });
      }
    }, 50);
  }

  closeDropdown() {
    if (this.isOpen) {
      this.isOpen = false;
      this.onTouch();
    }
  }

  selectOption(option: SelectOption) {
    this.value = option.value;
    this.selectedLabel = option.label;
    this.onChange(this.value);
    this.onSelect.emit(this.value);
    this.closeDropdown();
  }

  // Implementation of ControlValueAccessor
  writeValue(value: any): void {
    this.value = value;
    const option = this.options.find((opt) => opt.value === value);
    this.selectedLabel = option ? option.label : null;
  }

  registerOnChange(fn: any): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: any): void {
    this.onTouch = fn;
  }

  setDisabledState?(isDisabled: boolean): void {
    this.disabled = isDisabled;
  }
}
