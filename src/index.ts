/**
 * OpenAutonomyX SDK
 * The global creative platform - Create once, publish everywhere, in any language
 *
 * @package openautonomyx
 * @version 1.0.0
 */

// Core exports
export * from './types';
export * from './services/content';
export * from './services/publishing';
export * from './services/analytics';
export * from './services/team';
export * from './services/templates';
export * from './services/ai';
export * from './services/multilingual';

// SDK initialization
export { OpenAutonomyX } from './openautonomyx';

// Utilities
export * from './utils/errors';
export * from './utils/validators';

// Default export
import { OpenAutonomyX } from './openautonomyx';
export default OpenAutonomyX;
