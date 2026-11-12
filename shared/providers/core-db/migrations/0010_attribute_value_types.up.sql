-- Новые enum-значения должны быть зафиксированы до их использования в 0011.
ALTER TYPE attribute_value_type ADD VALUE IF NOT EXISTS 'decimal';
ALTER TYPE attribute_value_type ADD VALUE IF NOT EXISTS 'boolean';
